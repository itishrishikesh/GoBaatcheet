import React, { useState } from 'react'
import { UsernameHolder } from './component/UsernameHolder'
import { ChatContainer } from './component/ChatContainer'
import { User } from './models/User'
import { Message } from './models/Message'
import { Authorizer, AuthorizerProvider } from '@authorizerdev/authorizer-react'
import "./App.css"

function App() {
  const [username, setUsername] = useState<string>("")
  const [receiver, setReceiver] = useState<string>("")
  const [websocket, setWebsocket] = useState<WebSocket>()
  const [messages, setMessages] = useState<Message[]>([])

  const setupConnection = (response: any) => {
    console.log("authorizer response", response)
    const localUsername = response.user.email
    const ws = new WebSocket("ws://localhost:8080/ws?access_token=" + response.access_token);
    setUsername(localUsername)
    setWebsocket(ws)
    if(!ws) {
      console.error("E#1QGGEL - Websocket connection is null!")
      return;
    }
    ws.addEventListener("open", () => ws.send(JSON.stringify(new User(localUsername ?? ""))))
    ws.addEventListener("message", (event) => setMessages(prev => [...prev, JSON.parse(event.data)]))
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
    setMessages(prev => [...prev, msgToSend])
  }

  return (
    <>
      {
        username.length ?
          <div className="parent">
            <h4>GoBaatcheet</h4>
            <UsernameHolder sender={username} receiver={receiver} changeReceiver={changeReceiver} />
            <hr />
            <br />
            <ChatContainer send={send} messages={messages} />
          </div>
          :
          <div className="login-outer-container">
            <AuthorizerProvider
                config={{
                  authorizerURL: 'http://localhost:8082',
                  redirectURL: window.location.origin,
                  clientID: '7ea6b8a6-77f5-432e-b501-a51fe02c40c7'
                }}>

              <Authorizer onLogin={(response) => setupConnection(response)} onSignup={(response) => setupConnection(response)} />

            </AuthorizerProvider>
          </div>
      }
    </>
  )
}

export default App
