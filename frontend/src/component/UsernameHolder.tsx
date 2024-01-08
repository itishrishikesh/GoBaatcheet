export const UsernameHolder = (props: {sender: string, receiver: string}) => {
    return <div className="usernames">{props.sender} | {props.receiver}</div>
}