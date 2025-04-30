import { useState } from 'react'
import Axios, { type AxiosResponse } from 'axios'
import Button from '@mui/material/Button'

interface ResponseData {
  message: string
}

export default function GetRequestButton () {
  const [data, setData] = useState<ResponseData | null>(null)

  const handleClick = async (): Promise<void> => {
    try {
      const response: AxiosResponse<ResponseData> = await Axios.get('http://localhost:8080/hello')
      setData(response.data)
    } catch (error) {
      console.error(error)
    }
  }

  return (
        <div>
            <Button onClick={handleClick} variant="contained">Make GET request</Button>
            <div>{(data != null) ? data.message : 'No data to be displayed right now'}</div>
        </div>
  )
}
