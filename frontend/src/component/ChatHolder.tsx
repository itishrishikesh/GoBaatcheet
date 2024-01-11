import { Message } from '../models/Message'

export const  ChatHolder = (props: {messages: Message[]}) => {
    return <div className="chatholder">
        {
            props.messages.map(
                (value) => <div>{
                    MessageComponent(value)
                }</div>
            )
        }
    </div>;
}

const MessageComponent = (message: Message) => {
    return <div>
        <span className='message'>{message.message}</span>
        <span className='small'> from {message.sender} to {message.receiver}</span>
    </div>
}