import { useState } from "react";
import { useNavigate,Link} from "react-router-dom";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import api from "@/services/api";

export default function Login() {
  const navigate = useNavigate();

  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const handleLogin = async (e) => {
    e.preventDefault();

    try {
      const res = await api.post("/login", {
        username,
        password,
      });

      localStorage.setItem("token", res.data.token);

      alert("Login successful");

      navigate("/notes");
    } catch (err) {
      console.log(err);
      alert(err.response?.data || "Login failed");
    }
  };

  return (
    <>
    <div className="flex min-h-screen items-center justify-center">
      <Card className="w-[400px]">
        <CardHeader>
          <CardTitle>Login</CardTitle>
        </CardHeader>

        <CardContent>
          <form onSubmit={handleLogin} className="space-y-4">
            <Input
              placeholder="Username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
            />

            <Input
              type="password"
              placeholder="Password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />

            <Button className="w-full" type="submit">
              Login
            </Button>
            <p className="text-center text-sm">
              new user?{" "}
              <Link to="/signup" className="text-blue-500 hover:underline">
                SignUp
              </Link>
            </p>
          </form>
        </CardContent>
      </Card>
    </div>
    </>
  );
}
