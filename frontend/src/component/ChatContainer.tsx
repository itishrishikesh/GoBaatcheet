import { ChatHolder } from "./ChatHolder";
import { ChatSender } from "./ChatSender";

export const  ChatContainer = (props: {messages: string[], send: any}) => {
    return <div>
        <ChatHolder messages={props.messages}/>
        <ChatSender send={props.send}/>
    </div>;
}