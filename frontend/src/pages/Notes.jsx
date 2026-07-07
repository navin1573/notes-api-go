import { useEffect, useState } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import api from "@/services/api";

export default function Notes() {
  const [notes, setNotes] = useState([]);
  const [text, setText] = useState("");
  const [editingId, setEditingId] = useState(null);
const [editingText, setEditingText] = useState("");
const [originalText, setOriginalText] = useState("");

  const startEdit = (note) => {
  setEditingId(note.id);
  setEditingText(note.text);
  setOriginalText(note.text)
  };

  const saveEdit = async () => {
  await api.put(`/notes/${editingId}`, {
    text: editingText,
  });

  setEditingId(null);
  setEditingText("");

  fetchNotes();
  };
  const fetchNotes = async () => {
    try {
      const res = await api.get("/notes");
      setNotes(Array.isArray(res.data) ? res.data : []);
    } catch (err) {
      console.log(err);
    }
  };

  useEffect(() => {
    fetchNotes();
  }, []);

  const addNote = async () => {
    if (!text) return;

    await api.post("/notes", {
      text,
    });

    setText("");
    fetchNotes();
  };

  const deleteNote = async (id) => {
    await api.delete(`/notes/${id}`);
    fetchNotes();
  };

  const logout = () => {
    localStorage.removeItem("token");
    window.location.href = "/";
  };

  return (<>
  <h1 className="mt-10 text-5xl text-center font-bold onClick=">Quick Note</h1>
    <div className="max-w-2xl mx-auto mt-10 space-y-4">
      <div className="flex gap-2">
        <Input className="text-white text-lg font-semibold" 
          placeholder="Write a note..."
          value={text}
          onChange={(e) => setText(e.target.value)}
        />

        <Button className="px-5 py-5" onClick={addNote}>
          Add
        </Button>

        <Button className="px-5 py-5"variant="destructive" onClick={logout}>
          Logout
        </Button>
      </div>
      {(notes || []).map((note) => (
        <Card key={note.id}>
         <CardContent className="text-xl flex justify-between p-5">
  <div>
    {editingId === note.id ? (
  <Input className="text-black-500 text-2xl"
    value={editingText}
    onChange={(e) => setEditingText(e.target.value)}
  />
) : (
  <p>{note.text}</p>
)}
  </div>

  <div className="flex gap-2">
    {editingId === note.id ? (
  <Button
  disabled={editingText.trim() === originalText.trim()}
  onClick={saveEdit}
>
  Save
</Button>
) : (
  <Button onClick={() => startEdit(note)}>
    Edit
  </Button>
)}

    <Button
      variant="destructive"
      onClick={() => deleteNote(note.id)}
    >
      Delete
    </Button>
  </div>
</CardContent> 
        </Card>
      ))}
    </div>
    </>
  );
}