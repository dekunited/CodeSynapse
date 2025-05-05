import pandas as pd
import csv
import os
import requests
import re
import nltk
nltk.download('punkt_tab')
from nltk.tokenize import word_tokenize
from nltk.translate.bleu_score import sentence_bleu, SmoothingFunction
from codebleu import calc_codebleu
import time


lang_names = {
    "Python": "python",
    "Java": "java",
    "C++": "cpp"
}


# A simplified CodeBLEU implementation that focuses on basic aspects
def calculate_bleu(reference, candidate):
    """Calculate BLEU score between reference and candidate."""
    # Tokenize reference and candidate
    reference_tokens = nltk.word_tokenize(reference)
    candidate_tokens = nltk.word_tokenize(candidate)

    # Calculate BLEU score with smoothing
    smooth = SmoothingFunction().method1

    # Calculate with different n-gram weights
    weights_1 = (1.0, 0, 0, 0)  # BLEU-1
    weights_2 = (0.5, 0.5, 0, 0)  # BLEU-2
    weights_4 = (0.25, 0.25, 0.25, 0.25)  # BLEU-4

    try:
        bleu_1 = sentence_bleu([reference_tokens], candidate_tokens, weights=weights_1, smoothing_function=smooth)
        bleu_2 = sentence_bleu([reference_tokens], candidate_tokens, weights=weights_2, smoothing_function=smooth)
        bleu_4 = sentence_bleu([reference_tokens], candidate_tokens, weights=weights_4, smoothing_function=smooth)
        return bleu_1, bleu_2, bleu_4
    except Exception as e:
        print(f"Error calculating BLEU: {e}")
        return 0, 0, 0


def calculate_keyword_match(reference, candidate, language):
    """Calculate keyword matching score based on language keywords."""
    # Define keywords for each language
    keywords = {
        'python': ['def', 'class', 'if', 'elif', 'else', 'for', 'while', 'try', 'except',
                  'import', 'from', 'return', 'with', 'as', 'True', 'False', 'None',
                  'and', 'or', 'not', 'in', 'is', 'lambda'],
        'java': ['public', 'private', 'protected', 'class', 'interface', 'extends', 'implements',
                'void', 'int', 'boolean', 'String', 'if', 'else', 'for', 'while', 'try', 'catch',
                'return', 'new', 'this', 'super', 'static', 'final', 'null', 'true', 'false'],
        'cpp': ['int', 'float', 'double', 'char', 'bool', 'void', 'class', 'struct', 'enum',
               'const', 'static', 'if', 'else', 'for', 'while', 'switch', 'case', 'break',
               'return', 'try', 'catch', 'throw', 'new', 'delete', 'public', 'private', 'protected']
    }

    if language not in keywords:
        print(f"Warning: Keywords not defined for {language}")
        return 0

    # Tokenize reference and candidate
    ref_tokens = set(nltk.word_tokenize(reference))
    cand_tokens = set(nltk.word_tokenize(candidate))

    # Find language keywords in reference and candidate
    ref_keywords = set(ref_tokens.intersection(keywords[language]))
    cand_keywords = set(cand_tokens.intersection(keywords[language]))

    # Calculate F1 score for keyword matching
    if not ref_keywords:
        return 0

    precision = len(ref_keywords.intersection(cand_keywords)) / len(cand_keywords) if cand_keywords else 0
    recall = len(ref_keywords.intersection(cand_keywords)) / len(ref_keywords) if ref_keywords else 0

    if precision + recall == 0:
        return 0

    f1 = 2 * precision * recall / (precision + recall)
    return f1


def calculate_codebleu(reference, candidate, target_lang):
    """Calculate a simplified version of CodeBLEU."""
    # Get target language
    target_lang = lang_names[target_lang]

    # Calculate BLEU scores
    bleu_1, bleu_2, bleu_4 = calculate_bleu(reference, candidate)

    # Calculate keyword matching score
    keyword_match = calculate_keyword_match(reference, candidate, target_lang)

    # Weighted sum (simplified CodeBLEU)
    # 0.7 * BLEU + 0.3 * keyword_match
    codebleu = 0.7 * bleu_4 + 0.3 * keyword_match

    return {
        'bleu_1': bleu_1,
        'bleu_2': bleu_2,
        'bleu_4': bleu_4,
        'keyword_match': keyword_match,
        'codebleu': codebleu
    }


def perform_translations(code_df, start_lang, target_lang):
    # Iterate over the first 100 samples
    translation_info_list = []
    for index, row in code_df.head(100).iterrows():
        # Obtain the starting code and the ground truth translation
        start_lang_code = row[start_lang].rstrip('\n')
        target_lang_code_truth = row[target_lang].rstrip('\n')

        # Create the prompt to pass into the LLM
        prompt = f"Translate the following code from {start_lang} to {target_lang} without any comments in the translated code: {start_lang_code}"
        
        # Obtain the predicted translation
        max_retries = 10
        retries = 0
        while (retries < max_retries):
            start_time = time.time()
            response = requests.post(
                "http://localhost:11434/api/generate",
                json={
                    "model": "deepseek-coder:6.7b",
                    "prompt": prompt,
                    "stream": False
                }
            )
            inference_time = time.time() - start_time
            result = response.json()
            output = result.get("response", "").strip()

            # Retry if the code block is not present of 
            if (response.status_code != 200 or "```" not in output):
                retries += 1
                if retries < max_retries: 
                    continue
                else:
                    print(f"[ERROR] Failed to translate problem {row["problem_id"]} from {start_lang} to {target_lang}.")
                    break
            else:
                # Extract and clean the code block if present
                output = output.split("```")[1]
                output = " ".join(output.split("\n")[1:])
                output = re.sub(r"\s+", " ", output)
                if (target_lang == "Java"):
                    output = output.replace("\n", " ")
                else:
                    output = output.replace("\n", "NEW_LINE")
                target_lang_code_pred = output

                # Log current translation info
                translation_info = {}
                translation_info["example_id"] = row["problem_id"]
                translation_info["start_language"] = start_lang
                translation_info["target_language"] = target_lang
                translation_info["source_code"] = start_lang_code
                translation_info["reference_code"] = target_lang_code_truth
                translation_info["translated_code"] = target_lang_code_pred
                translation_info["inference_time"] = inference_time

                # Calculate bleu score
                metrics = calculate_codebleu(target_lang_code_truth, target_lang_code_pred, target_lang)
                translation_info.update(metrics)
                translation_info_list.append(translation_info)
                break
            
    # Remove old file if necessary
    filename = f"./data-processing/deepseek/{lang_names[start_lang]}_to_{lang_names[target_lang]}.csv"
    if os.path.exists(filename):
        os.remove(filename)
    
    # Write the translation info to a csv file
    with open(filename, 'w', newline='') as file:
        writer = csv.DictWriter(file, fieldnames=translation_info_list[0].keys())
        writer.writeheader()
        writer.writerows(translation_info_list)
    return

def main():
    # Define the languages that will be tested
    lang_combos = [
        ("Python", "C++"), 
        ("Python", "Java"),
        ("C++", "Python"),
        ("C++", "Java"),
        ("Java", "Python"),
        ("Java", "C++")]
    
    # Read in the csv with equivalent code in all languages
    code_df = pd.read_csv("all_languages_aligned.csv")

    # Obtain 100 translations for each language combo
    for lang_combo in lang_combos:
        perform_translations(code_df, lang_combo[0], lang_combo[1])


if __name__ == "__main__":
    main()