import { useState } from "react"

export const ChatSender = (props: { send: any }) => {
    const [message, setMessage] = useState("")

    return <div>
        <input className="send_input" name="message" onKeyDown={(event) => {
            if (event.key == "Enter") {
                props.send(message)
            }
        }} onChange={(e) => setMessage(e.target.value)}></input>
        <button onClick={() => props.send(message)}>Send</button>
    </div>
}