import { h } from "preact";
import { useContext, useState } from "preact/hooks";

import { AuthContext } from "../../utils/jwt";
import style from "./style.css";
import Icon from "../icon";
import Modal, { ModalContent } from "../modal";

const DeleteModal = ({ isOpen, close, filePath, refresh, setRefresh }) => {
  const { jwt } = useContext(AuthContext);
  const [error, setError] = useState("");

  if (!isOpen) {
    return null;
  }

  const onSubmit = (e) => {
    e.preventDefault();
    const sendDelete = async () => {
      const response = await fetch(`/api/delete${filePath}`, {
        method: "DELETE",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
          Authorization: `Bearer ${jwt}`,
        },
      });

      if (response.status !== 200) {
        setError(response.json()["err"]);
        return;
      }
      setError("");
      setRefresh(refresh + 1);
      close();
    };
    sendDelete();
  };

  return (
    <Modal isOpen={isOpen} close={close}>
      <div class={style.modalHeader}>
        Permanently delete file?
        <Icon name="close" onClick={close} title="Close" />
      </div>
      <ModalContent>
        <button type="submit" onClick={onSubmit} class={style.submit}>
          Delete
        </button>
        <button type="cancel" onClick={close} class={style.cancel}>
          Cancel
        </button>
        {error !== "" && <div class={style.error}>{error}</div>}
      </ModalContent>
    </Modal>
  );
};

const Delete = ({ filePath, refresh, setRefresh }) => {
  const [modalOpen, setModalOpen] = useState(false);

  return (
    <>
      <Icon name="remove" onClick={() => setModalOpen(true)} title="Delete" />
      <DeleteModal
        isOpen={modalOpen}
        close={() => setModalOpen(false)}
        filePath={filePath}
        refresh={refresh}
        setRefresh={setRefresh}
      />
    </>
  );
};

export default Delete;
