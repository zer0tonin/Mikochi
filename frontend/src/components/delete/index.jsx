import { h } from "preact";
import { useContext, useState } from "preact/hooks";

import { AuthContext } from "../../jwt";
import Icon from "../icon";
import Modal, { ModalContent } from "../modal";
import "./style.css";
import handleError from "../../error";
import Toast from "../toast";

const DeleteModal = ({ isOpen, close, filePath, refresh, setRefresh, setSuccess }) => {
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
        setError(await handleError(response));
        return;
      }
      setError("");
      setRefresh(refresh + 1);
      close();
      setSuccess(true);
      setTimeout(() => {
        setSuccess(false);
      }, 2000)
    };
    sendDelete();
  };

  return (
    <Modal isOpen={isOpen} close={close}>
      <div class="modalHeader">
        Permanently delete file?
        <Icon name="close" onClick={close} title="Close" />
      </div>
      <ModalContent>
        <button type="submit" onClick={onSubmit} class="submit">
          Delete
        </button>
        <button type="cancel" onClick={close} class="cancel">
          Cancel
        </button>
        {error !== "" && <div class="error">{error}</div>}
      </ModalContent>
    </Modal>
  );
};

const Delete = ({ filePath, refresh, setRefresh }) => {
  const [modalOpen, setModalOpen] = useState(false);
  const [success, setSuccess] = useState(false);

  return (
    <>
      <Icon name="remove" onClick={() => setModalOpen(true)} title="Delete" />
      <DeleteModal
        isOpen={modalOpen}
        close={() => setModalOpen(false)}
        filePath={filePath}
        refresh={refresh}
        setRefresh={setRefresh}
        setSuccess={setSuccess}
      />
      {success && <Toast text='File Deleted Successfully' isVisible={success} />}
    </>
  );
};

export default Delete;
