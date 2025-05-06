import { useState } from 'react'
import axios from 'axios'
import '../../styles/translation_ui.css'

interface TranslationRequest {
  translation: string;
  code: string;
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

      const response = await axios.post('http://localhost:8080/api/translate', request)

      setTranslatedCode(response.data.translatedCode)
      setModelUsed(response.data.modelUsed)
    } catch (err) {
      console.error('Translation error:', err)
      setError('Failed to translate code. Please try again.')
    } finally {
      setLoading(false)
    }
  }

  const languages = [
    { value: 'python', label: 'Python' },
    { value: 'java', label: 'Java' },
    { value: 'cpp', label: 'C++' },
  ]

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
            {languages.map(lang => (
              <option key={lang.value} value={lang.value}>{lang.label}</option>
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
            {languages.map(lang => (
              <option key={lang.value} value={lang.value}>{lang.label}</option>
            ))}
          </select>
        </div>
      </div>

      <div className="editor-row">
        <div className="editor-column">
          <label htmlFor="source-code">From</label>
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
          <label htmlFor="translated-code">To</label>
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
