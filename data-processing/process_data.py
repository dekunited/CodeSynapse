import os
import pandas as pd
from collections import defaultdict

# Base path to your XLCoST data
base_path = "./pair_data_tok_full"

# All languages to include
languages = ["C", "C++", "Java", "Python", "C#", "Javascript", "PHP"]

# Standardize language names
def standardize_lang(lang):
    lang_map = {
        "c": "C",
        "c++": "C++",
        "cpp": "C++",
        "java": "Java",
        "python": "Python",
        "py": "Python",
        "c#": "C#",
        "csharp": "C#",
        "javascript": "Javascript",
        "js": "Javascript",
        "php": "PHP"
    }
    return lang_map.get(lang.lower(), lang)

# Function to read map file - just gets the problem IDs in order
def read_map_file(file_path):
    if not os.path.exists(file_path):
        print(f"Warning: Map file not found: {file_path}")
        return []

    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            # Extract just the problem IDs, one per line
            problem_ids = [line.strip().split('-')[0] for line in f]
        return problem_ids
    except Exception as e:
        print(f"Error reading map file: {e}")
        return []

# Function to read code file - each line is a complete program
def read_code_file(file_path):
    if not os.path.exists(file_path):
        print(f"Warning: Code file not found: {file_path}")
        return []

    try:
        with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
            # Each line is a complete program
            return [line for line in f]
    except Exception as e:
        print(f"Error reading code file: {e}")
        return []


# Dictionary to store code for each language and problem ID
problem_code = {lang: {} for lang in languages}
all_problem_ids = set()

# Find all language pair directories
for dir_name in os.listdir(base_path):
    dir_path = os.path.join(base_path, dir_name)
    if not os.path.isdir(dir_path):
        continue

    # Parse language pair from directory name
    parts = dir_name.split('-')
    if len(parts) != 2:
        continue

    lang1, lang2 = parts
    lang1 = standardize_lang(lang1)
    lang2 = standardize_lang(lang2)

    # Skip if neither language is in our target list
    if lang1 not in languages and lang2 not in languages:
        continue

    print(f"Processing {lang1}-{lang2} directory")

    # Find map files
    lang1_map_path = None
    for f in os.listdir(dir_path):
        if f.startswith("train") and lang1 in f and "map" in f:
            lang1_map_path = os.path.join(dir_path, f)
            break

    lang2_map_path = None
    for f in os.listdir(dir_path):
        if f.startswith("train") and lang2 in f and "map" in f:
            lang2_map_path = os.path.join(dir_path, f)
            break

    if not lang1_map_path or not lang2_map_path:
        print(f"  Could not find map files for {dir_name}")
        continue

    # Find code files with appropriate extensions
    extensions = {
        "C": [".c"],
        "C++": [".cpp", ".c++"],
        "Java": [".java"],
        "Python": [".py", ".python"],
        "C#": [".cs", ".c#"],
        "Javascript": [".js", ".javascript"],
        "PHP": [".php"]
    }

    lang1_code_path = None
    for ext in extensions.get(lang1, []):
        for f in os.listdir(dir_path):
            if f.startswith("train") and f.endswith(ext):
                lang1_code_path = os.path.join(dir_path, f)
                break
        if lang1_code_path:
            break

    lang2_code_path = None
    for ext in extensions.get(lang2, []):
        for f in os.listdir(dir_path):
            if f.startswith("train") and f.endswith(ext):
                lang2_code_path = os.path.join(dir_path, f)
                break
        if lang2_code_path:
            break

    if not lang1_code_path or not lang2_code_path:
        print(f"  Could not find code files for {dir_name}")
        continue

    # Read map files to get problem IDs
    lang1_problem_ids = read_map_file(lang1_map_path)
    lang2_problem_ids = read_map_file(lang2_map_path)

    # Read code files
    lang1_code_lines = read_code_file(lang1_code_path)
    lang2_code_lines = read_code_file(lang2_code_path)

    # Check for mismatches
    if len(lang1_problem_ids) != len(lang1_code_lines):
        print(f"  Warning: Number of {lang1} problem IDs ({len(lang1_problem_ids)}) " +
              f"doesn't match number of code lines ({len(lang1_code_lines)})")

    if len(lang2_problem_ids) != len(lang2_code_lines):
        print(f"  Warning: Number of {lang2} problem IDs ({len(lang2_problem_ids)}) " +
              f"doesn't match number of code lines ({len(lang2_code_lines)})")

    # Map problem IDs to code
    for i, problem_id in enumerate(lang1_problem_ids):
        if i < len(lang1_code_lines):
            problem_code[lang1][problem_id] = lang1_code_lines[i]
            all_problem_ids.add(problem_id)

    for i, problem_id in enumerate(lang2_problem_ids):
        if i < len(lang2_code_lines):
            problem_code[lang2][problem_id] = lang2_code_lines[i]
            all_problem_ids.add(problem_id)

# Find problem IDs with implementations in all languages
common_problems = []
for problem_id in all_problem_ids:
    if all(problem_id in problem_code[lang] for lang in languages):
        common_problems.append(problem_id)

print(f"Found {len(common_problems)
               } problems with implementations in all languages")

# Create DataFrame with aligned programs
data = []
for problem_id in common_problems:
    row = {'problem_id': problem_id}
    for lang in languages:
        row[lang] = problem_code[lang][problem_id]
    data.append(row)

# Create and save the DataFrame
df = pd.DataFrame(data)
print(f"Created DataFrame with {len(df)} rows and {len(df.columns)} columns")
df.to_csv('all_languages_aligned.csv', index=False)
print("Saved CSV file: all_languages_aligned.csv")

# Create a summary of coverage
coverage = {lang: len(problem_code[lang]) for lang in languages}
print("\nProblem coverage by language:")
for lang, count in coverage.items():
    print(f"  {lang}: {count} problems")

# Print a sample row if available
if len(df) > 0:
    print("\nSample of first row (truncated):")
    sample_row = df.iloc[0].copy()
    for lang in languages:
        code_sample = sample_row[lang]
        sample_row[lang] = code_sample[:100] + \
            '...' if len(code_sample) > 100 else code_sample
    print(sample_row)
