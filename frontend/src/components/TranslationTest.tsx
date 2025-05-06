import { useState, useEffect } from 'react'
import axios from 'axios'
import '../../styles/translation_ui.css'

interface TranslationRequest {
  translation: string;
  code: string;
  model?: string;
}

interface TranslationResponse {
  translatedCode: string;
  modelUsed: string;
}

interface Language {
  value: string;
  label: string;
  benchmarked: boolean;
}

interface Model {
  value: string;
  label: string;
}

// Map of benchmarked language pairs to their best models (using API values)
const bestTranslationMap: Record<string, string> = {
  'java-python': 'llama-3.2-3b',
  'java-cpp': 'llama-3.2-3b',
  'python-java': 'llama-3.2-3b',
  'python-cpp': 'llama-3.2-3b',
  'cpp-python': 'llama-3.2-3b',
  'cpp-java': 'deepseek-6.7b',
  'go-python': 'gpt4o',
}

function decodePseudoCode(input: string): string {
  const cleanedInput = input
    .replace(/â–|␣/g, ' ') // Replace weird space chars
    .replace(/\bINDENT\b/g, '') // Remove INDENT for spacing
    .replace(/\bDEDENT\b/g, '') // Remove DEDENT

  const tokens = cleanedInput.split('NEW_LINE')
  const lines: string[] = []
  let indentLevel = 0
  const indent = () => '  '.repeat(indentLevel)

  for (let token of tokens) {
    token = token.trim()
    if (!token) continue

    if (token === 'INDENT') {
      indentLevel++
    } else if (token === 'DEDENT') {
      indentLevel = Math.max(0, indentLevel - 1)
    } else {
      const segments = token.split(';').map(s => s.trim()).filter(Boolean)
      for (let segment of segments) {
        lines.push(indent() + segment + (token.includes(';') ? ';' : ''))
      }
    }
  }

  return lines.join('\n')
}

export default function TranslationTest() {
  const [sourceLanguage, setSourceLanguage] = useState<string>('python')
  const [targetLanguage, setTargetLanguage] = useState<string>('java')
  const [sourceCode, setSourceCode] = useState<string>('')
  const [translatedCode, setTranslatedCode] = useState<string>('')
  const [selectedModel, setSelectedModel] = useState<string>('')
  const [defaultModel, setDefaultModel] = useState<string>('')
  const [isDefaultModel, setIsDefaultModel] = useState<boolean>(true)
  const [loading, setLoading] = useState<boolean>(false)
  const [error, setError] = useState<string | null>(null)

  // Define available models with consistent values for API
  const models: Model[] = [
    { value: 'llama-3.2-3b', label: 'Llama-3.2/3B' },
    { value: 'deepseek-6.7b', label: 'Deepseek-coder/6.7B' },
    { value: 'phi-2.7b', label: 'Phi/2.7B' },
    { value: 'gpt4o', label: 'GPT4o' }
  ];

  // Determine the best model based on language pair
  // Default to GPT4o for unbenchmarked pairs
  const getBestModel = (source: string, target: string): string => {
    const key = `${source}-${target}`

    if (key in bestTranslationMap) {
      return bestTranslationMap[key]
    }

    return 'gpt4o'
  }

  // Initialize default model on first load
  useEffect(() => {
    const best = getBestModel(sourceLanguage, targetLanguage)
    setDefaultModel(best)
    setSelectedModel(best)
  }, [])

  // Update the default model when source or target languages change
  useEffect(() => {
    const bestModel = getBestModel(sourceLanguage, targetLanguage)
    setDefaultModel(bestModel)
    setSelectedModel(bestModel)
    setIsDefaultModel(true)
  }, [sourceLanguage, targetLanguage])

  // Handle model selection changes
  const handleModelChange = (model: string) => {
    setSelectedModel(model)
    setIsDefaultModel(model === getBestModel(sourceLanguage, targetLanguage))
  }

  const handleSubmit = async (): Promise<void> => {
    setLoading(true)
    setError(null)

    try {
      // Apply the decodePseudoCode transformation only at submission time, so the user can still type
      const processedCode = decodePseudoCode(sourceCode)

      const request: TranslationRequest = {
        translation: `${sourceLanguage}-${targetLanguage}`,
        code: processedCode,
        model: selectedModel
      }

      const response = await axios.post<TranslationResponse>('http://localhost:8080/api/translate', request)
      const responseModel = response.data.modelUsed
      const knownModel = models.find(m => m.value === responseModel || m.label === responseModel)

      setTranslatedCode(response.data.translatedCode)
      if (knownModel) {
        setSelectedModel(knownModel.value)
      }

    } catch (err) {
      console.error('Translation error:', err)
      setError('Failed to translate code. Please try again.')
    } finally {
      setLoading(false)
    }
  }

  const languages: Language[] = [
    // Benchmarked languages
    { value: 'python', label: 'Python', benchmarked: true },
    { value: 'java', label: 'Java', benchmarked: true },
    { value: 'cpp', label: 'C++', benchmarked: true },

    // Unbenchmarked languages
    { value: 'csharp', label: 'C#', benchmarked: false },
    { value: 'javascript', label: 'JavaScript', benchmarked: false },
    { value: 'typescript', label: 'TypeScript', benchmarked: false },
    { value: 'go', label: 'Go', benchmarked: false },
    { value: 'rust', label: 'Rust', benchmarked: false }
  ]

  const renderLanguageOptions = () => {
    const benchmarkedLanguages = languages.filter(lang => lang.benchmarked)
    const unbenchmarkedLanguages = languages.filter(lang => !lang.benchmarked)

    return (
      <>
        <optgroup label="Benchmarked Languages">
          {benchmarkedLanguages.map(lang => (
            <option
              key={lang.value}
              value={lang.value}
            >
              {lang.label}
            </option>
          ))}
        </optgroup>
        <optgroup label="Unbenchmarked Languages">
          {unbenchmarkedLanguages.map(lang => (
            <option
              key={lang.value}
              value={lang.value}
              className="unbenchmarked-option"
            >
              {lang.label}
            </option>
          ))}
        </optgroup>
      </>
    )
  }

  return (
    <div className="container">
      <h1>CodeSynapse</h1>
      <h5>Bridge the Syntax, Power the Future!</h5>

      <div className="select-group">
        <div className="select-block align-left">
          <label htmlFor="source-lang">From</label>
          <select
            id="source-lang"
            value={sourceLanguage}
            onChange={(e) => setSourceLanguage(e.target.value)}
          >
            {renderLanguageOptions()}
          </select>
        </div>

        <div className="select-block align-center">
          <label htmlFor="llm-model">
            LLM Options
          </label>
          <select
            id="llm-model"
            value={selectedModel}
            onChange={(e) => handleModelChange(e.target.value)}
          >
            {models.map(model => (
              <option key={model.value} value={model.value}>
                {model.label}
                {model.value === defaultModel && ' (recommended)'}
              </option>
            ))}
          </select>
        </div>

        <div className="select-block align-right">
          <label htmlFor="target-lang">To</label>
          <select
            id="target-lang"
            value={targetLanguage}
            onChange={(e) => setTargetLanguage(e.target.value)}
          >
            {renderLanguageOptions()}
          </select>
        </div>
      </div>

      <div className="editor-row">
        <div className="editor-column">
          <textarea
            id="source-code"
            className="code-input"
            placeholder="Enter source code here..."
            value={sourceCode}
            onChange={(e) => setSourceCode(e.target.value)}
          ></textarea>
        </div>

        <div className="editor-column">
          <textarea
            id="translated-code"
            className="code-output"
            value={translatedCode}
            readOnly
          ></textarea>
        </div>
      </div>

      <button
        className="translation-button"
        onClick={handleSubmit}
        disabled={loading || !sourceCode.trim()}
      >
        {loading ? (
          <span className="loader"></span>
        ) : (
          `Translate Code`
        )}
      </button>

      {error && <p className="error-message">{error}</p>}
    </div>
  )
}
