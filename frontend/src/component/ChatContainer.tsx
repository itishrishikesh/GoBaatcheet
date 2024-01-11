import { ChatHolder } from "./ChatHolder";
import { ChatSender } from "./ChatSender";
import { Message } from '../models/Message'

export const  ChatContainer = (props: {messages: Message[], send: any}) => {
    return <div>
        <ChatHolder messages={props.messages}/>
        <ChatSender send={props.send}/>
    </div>;
}