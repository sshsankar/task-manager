import { useEffect, useState } from "react";
import { signInWithGoogle, auth } from "../lib/firebase";
import { useRouter } from "next/router";
import { onAuthStateChanged } from "firebase/auth";

export default function Login() {
  const router = useRouter();
  const [user, setUser] = useState(null);

  useEffect(() => {
    onAuthStateChanged(auth, (currentUser) => {
      setUser(currentUser);
      if (currentUser) {
        router.push("/");
      }
    });
  }, []);

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-100">
      <div className="bg-white shadow-lg rounded-lg p-8 max-w-md text-center">
        <h1 className="text-3xl font-bold mb-4">Welcome to Task Manager</h1>
        <p className="text-gray-600 mb-6">Sign in with Google to continue</p>
        <button onClick={signInWithGoogle} className="bg-blue-500 text-white px-6 py-2 rounded shadow hover:bg-blue-600 transition w-full">
          Sign in with Google
        </button>
      </div>
    </div>
  );
}
