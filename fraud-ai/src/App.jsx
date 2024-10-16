import { Fields } from "./Fields"
import "./App.css"
import { useState } from "react"
import axios from 'axios';

const SpamResultsFail = ({ reason, level }) => {
  return (
    <div className="h-[300px] w-[600px] border-4 border-red-500 rounded-md p-8">
      <p className="font-bold text-2xl pb-4">WARNING - Potential Spam</p>
      <p className="pb-2 mt-4 underline">Confidence Level:</p>
      <div className="bg-gray-200 rounded-full w-full mb-6">
        <div className="bg-red-500 text-white rounded-full p-2 flex justify-end" style={{ width: `${level}%`}}>
          <p className="font-bold text-sm">{level}%</p>
        </div>
      </div>
      <p className="pb-1 underline">Reason:</p>
      <p>{reason}</p>
    </div>
  )
}

const SpamResultsSuccess = ({ reason, level }) => {
  return (
    <div className="h-[300px] w-[600px] border-4 border-green-500 rounded-md p-8">
      <p className="font-bold text-2xl pb-4">Validated!</p>
      <p className="pb-2 mt-4 underline">Confidence Level:</p>
      <div className="bg-gray-200 rounded-full w-full mb-6">
        <div className="bg-green-500 text-white rounded-full p-2 flex justify-end" style={{ width: `${level}%`}}>
          <p className="font-bold text-sm">{level}%</p>
        </div>
      </div>
      <p className="pb-1 underline">Reason:</p>
      <p>{reason}</p>
    </div>
  )
}

function App() {
  const [isLoading, setIsLoading] = useState(false)
  const [resultsLoaded, setResultsLoaded] = useState(false)
  const [spamData, setSpamData] = useState({})


  const handleDataSubmission = async (data) => {
    try {
      const res = await axios.post("http://localhost:3001/api/fraud", data, {
        withCredentials: true
      })
      setResultsLoaded(true)
      setIsLoading(false)
      setSpamData(res.data)
    } catch (err) {
      setResultsLoaded(true)
      setIsLoading(false)
      console.log(err)
    }
  }

  return (
    <div className="h-screen w-screen flex overflow-hidden">
      <div className="bg-[#15246D] flex flex-col w-[300px] px-4 items-center justify-center text-white font-bold">
        <p className="text-xl">Justworks </p>
        <p className="text-4xl underline decoration-8 decoration-[#52B0FF] underline-offset-2">Hackday-8</p>
        <p className="mt-4 font-light">AI Fraud Detection</p>
        <p className="font-bold text-sm text-yellow-300">DEMO</p>
      </div>
      <div className="flex-auto p-8 m-auto">
        {
          resultsLoaded
          ? resultsLoaded.isSpam ? <SpamResultsFail reason={spamData.reason} level={spamData.level} /> : <SpamResultsSuccess reason={spamData.reason} level={spamData.level} />
          : <Fields isLoading={isLoading} setIsLoading={setIsLoading} handleDataSubmission={handleDataSubmission}/>
        }
      </div>
    </div>
  )
}

export default App
