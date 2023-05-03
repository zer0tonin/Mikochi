import { h } from "preact";
import { useContext, useState } from "preact/hooks";

import { AuthContext } from "../../utils/jwt";

const Login = () => {
    const {setJWT} = useContext(AuthContext)
    const [username, setUsername] = useState("")
    const [password, setPassword] = useState("")
    const [error, setError] = useState("")

    const onSubmit = (e) => {
        e.preventDefault()
        const postLogin = async () => {
            const response = await fetch('/api/login', {
                method: 'POST',
                headers: {
                  'Accept': 'application/json',
                  'Content-Type': 'application/json'
                },
                body: JSON.stringify({ username, password })
            })
            const json = response.json()
            if (response.status !== 200) {
                setError(json["err"])
                return
            }
            setJWT(json["token"])
        }
        postLogin()
    }

    return (
        <main>
            <form onSubmit={onSubmit}>
                <input
                    type="text"
                    value={username}
                    onInput={(e) => setUsername(e.target.value)}
                />
                <input
                    type="password"
                    value={password}
                    onInput={(e) => setPassword(e.target.value)}
                />
                <button type="submit">Login</button>
            </form>
        </main>
    )
}

export default Login
