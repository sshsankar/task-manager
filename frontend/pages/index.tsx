import { useEffect, useState } from "react";
import axios from "axios";

export default function Home() {
  const [tasks, setTasks] = useState([]);
  const [aiSuggestions, setAiSuggestions] = useState([]);
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");

  useEffect(() => {
    axios.get("http://localhost:8080/tasks").then((res) => {
      setTasks(res.data);
    });

    const socket = new WebSocket("ws://localhost:8080/ws");
    socket.onmessage = () => {
      axios.get("http://localhost:8080/tasks").then((res) => {
        setTasks(res.data);
      });
    };

    return () => socket.close();
  }, []);

  const addTask = async () => {
    try {
      console.log("Sending task:", { title, description }); // Debugging log

      const response = await axios.post("http://localhost:8080/tasks", { title, description });

      console.log("Task added successfully:", response.data); // Debugging log

      setTitle("");
      setDescription("");

      axios.get("http://localhost:8080/tasks").then((res) => {
        setTasks(res.data);
      });
    } catch (error) {
      console.error("Error adding task:", error);
      alert("Failed to add task. Check console for details.");
    }
  };

  // DELETE TASK FUNCTION
  const deleteTask = async (id: number) => {
    try {
      await axios.delete(`http://localhost:8080/tasks/${id}`);

      // Remove the deleted task from the UI
      setTasks(tasks.filter(task => task.id !== id)); 
    } catch (error) {
      console.error("Error deleting task:", error);
    }
  };

  const fetchTaskSuggestions = async () => {
    try {
      const response = await axios.get("http://localhost:8080/ai-tasks");

      if (!response.data.choices || response.data.choices.length === 0) {
        throw new Error("Invalid AI response structure");
      }

      const aiText = response.data.choices[0]?.message?.content;
      if (!aiText) {
        throw new Error("No content in AI response.");
      }

      const suggestions = aiText.split("\n").filter(task => task.trim() !== "");
      setAiSuggestions(suggestions);
    } catch (error) {
      console.error("Error fetching AI task suggestions:", error);
      setAiSuggestions(["Failed to get AI suggestions. Please try again."]);
    }
  };

  return (
    <div className="min-h-screen bg-gray-100 p-8">
      <div className="max-w-3xl mx-auto bg-white shadow-md rounded-lg p-6">
        <h1 className="text-3xl font-bold text-center mb-6">Task Dashboard</h1>

        {/* Add New Task Form */}
        <div className="mb-6 flex flex-col gap-4">
          <input
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="Task Title"
            className="border p-2 rounded w-full"
          />
          <input
            type="text"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="Task Description"
            className="border p-2 rounded w-full"
          />
          <button onClick={addTask} className="bg-green-500 text-white px-6 py-2 rounded shadow hover:bg-green-600 transition">
            Add Task
          </button>
        </div>

        {/* Task List */}
        <div className="mt-4">
          {tasks.length === 0 ? (
            <p className="text-center text-gray-500">No tasks available</p>
          ) : (
            tasks.map((task) => (
              <div key={task.id} className="border p-4 mb-4 rounded-lg shadow flex justify-between items-center">
                <div>
                  <h2 className="text-lg font-semibold">{task.title}</h2>
                  <p className="text-gray-600">{task.description}</p>
                  <p className="text-gray-500"><strong>Status:</strong> {task.status}</p>
                </div>
                <button
                  onClick={() => deleteTask(task.id)}
                  className="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600 transition"
                >
                  Delete
                </button>
              </div>
            ))
          )}
        </div>

        {/* AI Task Suggestions */}
        <button onClick={fetchTaskSuggestions} className="mt-6 bg-blue-500 text-white px-6 py-2 rounded shadow hover:bg-blue-600 transition w-full">
          Get AI Task Suggestions
        </button>

        {aiSuggestions.length > 0 && (
          <div className="mt-6 p-4 bg-gray-50 rounded-lg shadow">
            <h2 className="text-lg font-semibold">AI Suggested Tasks:</h2>
            <ul className="list-disc pl-6 text-gray-700">
              {aiSuggestions.map((suggestion, index) => (
                <li key={index}>{suggestion}</li>
              ))}
            </ul>
          </div>
        )}
      </div>
    </div>
  );
}
