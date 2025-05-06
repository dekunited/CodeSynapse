# Data processing
This directory contains the code + resulting CSV for parsing the full programs contained in the [XLCoST](https://github.com/reddy-lab-code-research/XLCoST) dataset. The goal was to gather a set of programs that were implemented in multiple languages in order to evaluate different LLMs for code-translation tasks.

Additionally, it also contains all code used for training and gathering our data/insights as to which models perform best for which language translation pairs. For more information about how the models performed, please visit their corresponding directories or view our final paper/presentation. 


## Dataset Info 
The *XLCoST* dataset contains mapping files to assist with finding which lines in each code file correspond to which specific problem. This differs between the snippets and full program data so *please note: we used the full program data and as such, this script is mainly for the full program data. It will need to be modified further to work on the snippets data*. 

In the case of the full program data, a single line in a code file is actually the entire program. The mapping file lines up directly with this, specifying which problem this program is for. For example:

```
train-C-map.jsonl
10113
....

// C code for the C-C# translation
train-C-C#-tok.c 
#include <stdio.h> NEW_LINE void count_setbit ( int N ) { int result = 0 ; for ( int i = 0 ; i < 32 ; i ++ ) { if ( ( 1 << i ) & N ) { result ++ ; } } printf ( " % d STRNEWLINE " , result ) ; } int main ( ) { int N = 43 ; count_setbit ( N ) ; return 0 ; }

This single line containing the entire program corresponds to problem 10113
```

## Dataset Gathering Script Info 
This script, again, is primarily for the program data. It will go through all translations, match the corresponding problem id to theline containing the program, and match them across languages. If a problem was implemented in all 7 languages, it was then added to our CSV. 

The script assumes you have the `pair_data_tok_full/` folder in your current directory.

The script can be modified to also include data if implemented for a specified minimum of languages (rather then all 7). So, if you wanted programs that were modified in at least 5 of the languages, you can modify this piece of the code to only require your specified amount:

```Python

# Find problem IDs with implementations in all languages
common_problems = []
for problem_id in all_problem_ids:
    if all(problem_id in problem_code[lang] for lang in languages):
        common_problems.append(problem_id)


```
