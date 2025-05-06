import { useState } from 'react'
import axios from 'axios'
import '../../styles/translation_ui.css'

interface TranslationRequest {
  translation: string;
  code: string;
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

function decodePseudoCode(input: string): string {
  const cleanedInput = input
    .replace(/â–|␣/g, ' ') // Replace weird space chars
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
  const [modelUsed, setModelUsed] = useState<string>('')
  const [loading, setLoading] = useState<boolean>(false)
  const [error, setError] = useState<string | null>(null)

  const handleSubmit = async (): Promise<void> => {
    setLoading(true)
    setError(null)

    try {
      const request: TranslationRequest = {
        translation: `${sourceLanguage}-${targetLanguage}`,
        code: sourceCode
      }

      const response = await axios.post<TranslationResponse>('http://localhost:8080/api/translate', request)

      setTranslatedCode(response.data.translatedCode)
      setModelUsed(response.data.modelUsed)
    } catch (err) {
      console.error('Translation error:', err)
      setError('Failed to translate code. Please try again.')
    } finally {
      setLoading(false)
    }
  }

  // Updated languages array with benchmarked property
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

  const llmModelOption = () => {
    const models = [
      { value: 'llama-3.2-3b', label: 'Llama-3.2/3B' },
      { value: 'deepseek-6.7b', label: 'Deepseek-coder/6.7B' },
      { value: 'phi-2.7b', label: 'Phi/2.7B' }
    ];
  
    return (
      <>
        {models.map(model => (
          <option key={model.value} value={model.value}>
            {model.label}
          </option>
        ))}
      </>
    );
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
          <label htmlFor="llm-model">LLM Model</label>
          <select
            id="llm-model"
            value={modelUsed}
            onChange={(e) => setModelUsed(e.target.value)}
          >
            {llmModelOption()}
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
            onChange={(e) => {
              const raw = e.target.value
              const formatted = decodePseudoCode(raw)
              setSourceCode(formatted)
            }}
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
        {loading ? <span className="loader"></span> : 'Translate Code'}
      </button>

      {error && <p className="error-message">{error}</p>}
    </div>
  )
}
