import React, { useState } from 'react'
import { UsernameHolder } from './component/UsernameHolder'
import { ChatContainer } from './component/ChatContainer'
import { User } from './models/User'
import { Message } from './models/Message'

function App() {
  const [username, setUsername] = useState<string>("")
  const [receiver, setReceiver] = useState<string>("")
  const [websocket, setWebsocket] = useState<WebSocket>()
  const [messages, setMessages] = useState<string[]>([])

  const logIn = () => {
    const localUsername: string | null = prompt("enter username")
    const ws = new WebSocket("ws://localhost:8080/ws")
    setUsername(localUsername ?? "")
    setWebsocket(ws)
    if(!ws) {
      console.error("E#1QGGEL - Websocket connection is null!")
      return;
    }
    ws.addEventListener("open", () => ws.send(JSON.stringify(new User(localUsername ?? ""))))
    ws.addEventListener("message", (event) => setMessages(prev => [...prev, JSON.parse(event.data).message]))
  }

  const changeReceiver = (e: React.ChangeEvent<HTMLInputElement>) => {
    setReceiver(e.target.value)
  }

  const send = (message: string) => {
    if(!websocket){
      console.error("E#1QGGKJ - Websocket connection is null!")
      return
    }
    const msgToSend = new Message(message, username, receiver)
    websocket.send(JSON.stringify(msgToSend))
    setMessages(prev => [...prev, message])
  }

  return (
    <>
      {
        username.length ?
          <div>
            <h4>GoBaatcheet</h4>
            <UsernameHolder sender={username} receiver={receiver} changeReceiver={changeReceiver} />
            <hr />
            <br />
            <ChatContainer send={send} messages={messages} />
          </div>
          :
          <div>Please <button onClick={() => logIn()}>enter</button> username!</div>
      }
    </>
  )
}

export default App
