import { h } from "preact";
import { useContext, useState } from "preact/hooks";

import { AuthContext } from "../../jwt";
import Icon from "../icon";
import Modal, { ModalContent } from "../modal";
import "./style.css";
import handleError from "../../error";
import Toast from "../toast";
import { signal } from "@preact/signals";
import { refresh } from "../../pages/Directory";

export const deleteFilePath = signal(null);

const DeleteModal = ({ isOpen, close, filePath, setSuccess }) => {
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
      refresh.value = refresh.value + 1;
      close();
      setSuccess(true);
      setTimeout(() => {
        setSuccess(false);
      }, 2000);
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
        <button type="submit" onClick={onSubmit} class="delete-submit">
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

const Delete = () => {
  const [success, setSuccess] = useState(false);

  return (
    <>
      <DeleteModal
        isOpen={deleteFilePath.value != null}
        close={() => (deleteFilePath.value = null)}
        filePath={deleteFilePath.value}
        setSuccess={setSuccess}
      />
      {success && (
        <Toast text="File Deleted Successfully" isVisible={success} />
      )}
    </>
  );
};

export default Delete;
