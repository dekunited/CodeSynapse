from flask import Flask, request, jsonify
from flask_cors import CORS
from transformers import AutoModelForSeq2SeqLM, AutoTokenizer
import torch

app = Flask(__name__)
CORS(app)

model_name = "Salesforce/codet5-base"
tokenizer = AutoTokenizer.from_pretrained(model_name)
model = AutoModelForSeq2SeqLM.from_pretrained(model_name)


@app.route('/translate', methods=['POST'])
def translate_code():
    data = request.get_json()
    prompt = data.get("prompt", "")

    if not prompt:
        return jsonify({"error": "No prompt provided"}), 400

    inputs = tokenizer(prompt, return_tensors="pt",
                       padding=True, truncation=True)

    with torch.no_grad():
        outputs = model.generate(**inputs, max_length=256)

    # generated = tokenizer.decode(outputs[0], skip_special_tokens=True)
    generated_code = tokenizer.decode(
        outputs[0].tolist(), skip_special_tokens=True)

    print("Code: ", generated_code)

    return jsonify({
        "translatedCode": generated_code,
        "modelUsed": model_name
    })


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=6969)
