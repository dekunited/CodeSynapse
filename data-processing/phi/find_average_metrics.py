import pandas as pd
import numpy as np
import csv
import os


def main():
    # List of csv files to examine
    files = [
        "./data-processing/phi/cpp_to_java.csv",
        "./data-processing/phi/cpp_to_python.csv",
        "./data-processing/phi/java_to_cpp.csv",
        "./data-processing/phi/java_to_python.csv",
        "./data-processing/phi/python_to_cpp.csv",
        "./data-processing/phi/python_to_java.csv"
    ]

    # List of metrics in each file that we want to find average of
    metrics = [
        "inference_time",
        "bleu_1",
        "bleu_2",
        "bleu_4",
        "keyword_match",
        "codebleu"
    ]

    # Dictionary key names for object storing the average metrics
    keys = [
        "avg_inference_time",
        "avg_bleu1_score",
        "avg_bleu2_score",
        "avg_bleu4_score",
        "avg_keyword_match",
        "avg_codebleu_score"
    ]

    # Calculate the average metric for all files
    avg_metrics_list = []
    for file in files:
        df = pd.read_csv(file)
        avg_metrics = {}
        avg_metrics["translation"] = file.split("/")[-1].split(".")[0]
        for i in range(len(metrics)):
            data = df[metrics[i]]
            avg = np.mean(data)
            avg_metrics[keys[i]] = avg
        avg_metrics["num_samples"] = 100
        avg_metrics_list.append(avg_metrics)

    # Remove old file if necessary
    filename = f"./data-processing/phi/phi_avg_metrics.csv"
    if os.path.exists(filename):
        os.remove(filename)
    
    # Write the translation info to a csv file
    with open(filename, 'w', newline='') as file:
        writer = csv.DictWriter(file, fieldnames=avg_metrics_list[0].keys())
        writer.writeheader()
        writer.writerows(avg_metrics_list)
    return


if __name__ == "__main__":
    main()