export const UsernameHolder = (props: {sender: string, receiver: string, changeReceiver: any}) => {

    return <div className="usernames">
            {props.sender} | <input onChange={props.changeReceiver} value={props.receiver} />
        </div>;
}