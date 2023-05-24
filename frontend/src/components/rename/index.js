import { h } from "preact";
import { useContext, useEffect, useRef, useState } from "preact/hooks";

import { AuthContext } from "../../utils/jwt";
import style from "./style.css";
import Icon from "../icon";

const RenameModal = ({ isOpen, close, filePath, refresh }) => {
  if (!isOpen) {
    return null;
  }

  const { jwt } = useContext(AuthContext);
  const [error, setError] = useState("");
  const [path, setPath] = useState(filePath);
  const ref = useRef();
  useEffect(() => {
    const handleClickOutside = (e) => {
      if (ref.current && !ref.current.contains(e.target)) {
        close();
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [ref]);

  const onSubmit = (e) => {
    e.preventDefault();
    const putMove = async () => {
      const response = await fetch(`/api/move${filePath}`, {
        method: "PUT",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
          Authorization: `Bearer ${jwt}`,
        },
        body: JSON.stringify({ newPath: path }),
      });

      if (response.status !== 200) {
        setError(json["err"]);
        return;
      }
      setError("");
      close();
      refresh();
    };
    putMove();
  };

  return (
    <div class={style.modal} ref={ref}>
      <div class={style.modalHeader}>
        Rename / Move
        <Icon name="close" onClick={close} title="Close" />
      </div>
      <div class={style.modalContent}>
        <form onSubmit={onSubmit}>
          <input
            type="text"
            value={path}
            class={style.input}
            onChange={(e) => setPath(e.target.value)}
          />
          <button type="submit" class={style.submit}>
            Rename
          </button>
          {error !== "" && <div class={style.error}>{error}</div>}
        </form>
      </div>
    </div>
  );
};

const Rename = ({ filePath, refresh }) => {
  const [modalOpen, setModalOpen] = useState(false);

  return (
    <>
      <Icon name="rename" onClick={() => setModalOpen(true)} title="Rename" />
      <RenameModal
        isOpen={modalOpen}
        close={() => setModalOpen(false)}
        filePath={filePath}
        refresh={refresh}
      />
    </>
  );
};

export default Rename;
