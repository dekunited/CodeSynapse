import { useState } from 'react'
import axios from 'axios'
import { Button, TextField, MenuItem, Select, FormControl, InputLabel, Box, Typography, CircularProgress, Paper } from '@mui/material'

interface TranslationRequest {
  translation: string;
  code: string;
}

/*
interface TranslationResponse {
  translatedCode: string;
  modelUsed: string;
}*/

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
    { value: 'javascript', label: 'JavaScript' },
    { value: 'go', label: 'Go' },
    { value: 'cpp', label: 'C++' },
    { value: 'csharp', label: 'C#' },
    { value: 'typescript', label: 'TypeScript' },
    { value: 'rust', label: 'Rust' }
  ]

  return (
    <Paper elevation={3} sx={{ p: 3, m: 2 }}>
      <Typography variant="h5" component="h2" gutterBottom>
        Code Translation
      </Typography>
      
      <Box sx={{ display: 'flex', gap: 2, mb: 2 }}>
        <FormControl sx={{ minWidth: 120 }}>
          <InputLabel id="source-language-label">From</InputLabel>
          <Select
            labelId="source-language-label"
            value={sourceLanguage}
            label="From"
            onChange={(e) => setSourceLanguage(e.target.value)}
          >
            {languages.map((lang) => (
              <MenuItem key={lang.value} value={lang.value}>{lang.label}</MenuItem>
            ))}
          </Select>
        </FormControl>
        
        <FormControl sx={{ minWidth: 120 }}>
          <InputLabel id="target-language-label">To</InputLabel>
          <Select
            labelId="target-language-label"
            value={targetLanguage}
            label="To"
            onChange={(e) => setTargetLanguage(e.target.value)}
          >
            {languages.map((lang) => (
              <MenuItem key={lang.value} value={lang.value}>{lang.label}</MenuItem>
            ))}
          </Select>
        </FormControl>
      </Box>
      
      <TextField
        label="Source Code"
        multiline
        rows={8}
        fullWidth
        value={sourceCode}
        onChange={(e) => setSourceCode(e.target.value)}
        variant="outlined"
        sx={{ mb: 2 }}
      />
      
      <Button 
        variant="contained" 
        onClick={handleSubmit}
        disabled={loading || !sourceCode.trim()}
        sx={{ mb: 2 }}
      >
        {loading ? <CircularProgress size={24} /> : 'Translate Code'}
      </Button>
      
      {error && (
        <Typography color="error" sx={{ mb: 2 }}>
          {error}
        </Typography>
      )}
      
      {translatedCode && (
        <Box>
          <Typography variant="subtitle2" sx={{ mb: 1 }}>
            Translated with {modelUsed}
          </Typography>
          <TextField
            label="Translated Code"
            multiline
            rows={8}
            fullWidth
            value={translatedCode}
            InputProps={{ readOnly: true }}
            variant="outlined"
          />
        </Box>
      )}
    </Paper>
  )
}
