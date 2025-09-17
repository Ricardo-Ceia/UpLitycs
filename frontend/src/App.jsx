import { useState } from 'react'
import Home from './Home'
import './App.css'

const user = {
  name: "Ricardo",
  city: "Lisbon"
}

function App() {
  return( 
    <div>
      <Home/>   
    </div>
  )
}

export default App
