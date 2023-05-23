import { h } from "preact";
import { useContext, useEffect, useRef, useState } from "preact/hooks";

import { AuthContext } from "../../utils/jwt";
import style from "./style.css";
import Icon from "../icon";


const DeleteModal = ({ isOpen, close, filePath, refresh }) => {
  if (!isOpen) {
    return null
  }

  const { jwt } = useContext(AuthContext);
  const [error, setError] = useState("");
  const ref = useRef();
  useEffect(() => {
    const handleClickOutside = (e) => {
      if (ref.current && !ref.current.contains(e.target)) {
        close()
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [ref]);

  const onSubmit = (e) => {
    e.preventDefault();
    const sendDelete = async () => {
      const response = await fetch(
        `/api/delete${filePath}`,
        {
          method: "DELETE",
          headers: {
            Accept: "application/json",
            "Content-Type": "application/json",
            Authorization: `Bearer ${jwt}`,
          },
        },
      );

      if (response.status !== 200) {
        setError(json["err"]);
        return;
      }
      setError("")
      close()
      refresh()
    }
    sendDelete();
  }

  return (
    <div class={style.modal} ref={ref}>
      <div class={style.modalHeader}>
        Permanently delete file?
        <Icon name="close" onClick={close} />
      </div>
      <div class={style.modalContent}>
        <button
          type="submit"
          onClick={onSubmit}
          class={style.submit}
        >
          Delete
        </button>
        <button
          type="cancel"
          onClick={close}
          class={style.cancel}
        >
          Cancel
        </button>
        {error !== "" && <div class={style.error}>{error}</div>}
      </div>
    </div>
  );
};

const Delete = ({ filePath, refresh }) => {
  const [modalOpen, setModalOpen] = useState(false);

  return (
    <>
      <Icon name="remove" onClick={() => setModalOpen(true)} />
      <DeleteModal
        isOpen={modalOpen}
        close={() => setModalOpen(false)}
        filePath={filePath}
        refresh={refresh}
      />
    </>
  );
};

export default Delete;