export const UsernameHolder = (props: {sender: string, receiver: string, changeReceiver: any}) => {

    return <div className="usernames">
           <span>{props.sender}</span> | <input className="receiverInput" placeholder="Enter receiver username" onChange={props.changeReceiver} value={props.receiver} />
        </div>;
}