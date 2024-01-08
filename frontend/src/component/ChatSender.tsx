import { useState } from "react"

export const ChatSender = (props: {send: any}) => {
    const [message, setMessage] = useState("")
    return <div>
        <input name="message" onChange={(e) => setMessage(e.target.value)}></input>
        <button onClick={() => props.send(message)}>Send</button>
    </div>
}