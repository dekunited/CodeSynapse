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

export default function TranslationTest() {
  const [sourceLanguage, setSourceLanguage] = useState<string>('')
  const [targetLanguage, setTargetLanguage] = useState<string>('')
  const [sourceCode, setSourceCode] = useState<string>('')
  const [translatedCode, setTranslatedCode] = useState<string>('')
  const [modelUsed, setModelUsed] = useState<string>('')
  const [loading, setLoading] = useState<boolean>(false)
  const [error, setError] = useState<string | null>(null)

  function decodePseudoCode(input: string): string {
      // const lines: string[] = []
      // const tokens = input.replace(/â–|␣/g, ' ').split('NEW_LINE')
      const cleanedInput = input.replace(/â–|␣/g, ' ')           // Replace weird space chars
        .replace(/\bINDENT\b/g, '  ')       // Remove INDENT
        .replace(/\bDEDENT\b/g, ' ')       // Remove DEDENT
      let indentLevel = 0
      const indent = () => '  '.repeat(indentLevel)

      const tokens = cleanedInput.split('NEW_LINE')
        const lines: string[] = []

        for (let token of tokens) {
          token = token.trim()
          if (!token) continue

          const segments = token.split(';').map(s => s.trim()).filter(Boolean)
          for (let segment of segments) {
            lines.push(segment)
          }
        }
            
      for (let token of tokens) {
        token = token.trim()
    
        if (token === 'INDENT') {
          indentLevel++
        } else if (token === 'DEDENT') {
          input.replace('DEDENT', ' ')
          indentLevel = Math.max(0, indentLevel - 1)
        } else if (token) {
          // Split at semicolons to create multiple lines from one token
          const segments = token.split(';').map(s => s.trim()).filter(Boolean)
          for (let segment of segments) {
            lines.push(indent() + segment)
          }
        }
      }
    
      return lines.join('\n')
    }
  
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

  const languages = [
    { value: 'python', label: 'Python' },
    { value: 'java', label: 'Java' },
    // { value: 'javascript', label: 'JavaScript' },
    // { value: 'go', label: 'Go' },
    { value: 'cpp', label: 'C++' },
    // { value: 'csharp', label: 'C#' },
    // { value: 'typescript', label: 'TypeScript' },
    // { value: 'rust', label: 'Rust' }
  ]

  return (
    <div className="container">
      <h1>CodeSynapse</h1>
      <h5>Bridge the Syntax, Power the Future!</h5>

      <div className="select-group">
        <div className="select-block">
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

        <div className="select-block">
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

      <textarea
        className="code-input"
        placeholder="Enter source code here..."
        rows={10}
        value={sourceCode}
        onChange={(e) => {
          const raw = e.target.value
          const formatted = decodePseudoCode(raw)
          setSourceCode(formatted)
        }}
      ></textarea>


      <button
        className="translation-button"
        onClick={handleSubmit}
        disabled={loading || !sourceCode.trim()}
      >
        {loading ? 'Translating...' : 'Translate Code'}
      </button>

      {error && <p className="error-message">{error}</p>}

      {translatedCode && (
        <div className="result-block">
          <p className="model-used">Translated with {modelUsed}</p>
          <textarea
            className="code-output"
            rows={10}
            value={translatedCode}
            readOnly
          ></textarea>
        </div>
      )}
    </div>
  )
}