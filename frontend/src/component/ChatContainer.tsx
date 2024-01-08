import React, { useState } from "react";
import { ChatHolder } from "./ChatHolder";
import { ChatSender } from "./ChatSender";

export const  ChatContainer = () => {
    const [messages, setMessages]  = useState<string[]>([])

    const send = (message: string) => {
        setMessages(prev => [...prev, message])
    }

    return <div>
        <ChatHolder messages={messages}/>
        <ChatSender send={send}/>
    </div>;
}