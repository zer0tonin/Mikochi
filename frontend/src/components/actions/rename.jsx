import { h } from "preact";
import { useContext, useEffect, useState } from "preact/hooks";

import { AuthContext } from "../../jwt";
import Icon from "../icon";
import "./style.css";

import Modal, { ModalContent, ModalHeader } from "../modal";
import handleError from "../../error";
import Toast from "../toast";

const RenameModal = ({
  isOpen,
  close,
  filePath,
  refresh,
  setRefresh,
  setSuccess,
}) => {
  const { jwt } = useContext(AuthContext);
  const [error, setError] = useState("");
  const [path, setPath] = useState(filePath);

  useEffect(() => {
    setPath(filePath);
  }, [filePath]);

  if (!isOpen) {
    return null;
  }

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
        setError(await handleError(response));
        return;
      }
      setError("");
      setRefresh(refresh + 1);
      close();
      setSuccess(true);
      setTimeout(() => {
        setSuccess(false);
      }, 2000);
    };
    putMove();
  };

  return (
    <Modal isOpen={isOpen} close={close}>
      <ModalHeader close={close}>Rename / Move</ModalHeader>
      <ModalContent>
        <form onSubmit={onSubmit}>
          <input
            type="text"
            value={path}
            class="renameInput"
            onChange={(e) => setPath(e.target.value)}
          />
          <button type="submit" class="submit">
            Rename
          </button>
          {error !== "" && <div class="error">{error}</div>}
        </form>
      </ModalContent>
    </Modal>
  );
};

export const RenameIcon = ({ filePath, refresh, setRefresh }) => {
  const [modalOpen, setModalOpen] = useState(false);
  const [success, setSuccess] = useState(false);

  return (
    <>
      <Icon name="rename" onClick={() => setModalOpen(true)} title="Rename" />
      <RenameModal
        isOpen={modalOpen}
        close={() => setModalOpen(false)}
        filePath={filePath}
        refresh={refresh}
        setRefresh={setRefresh}
        setSuccess={setSuccess}
      />
      {success && (
        <Toast text="File Renamed Successfully" isVisible={success} />
      )}
    </>
  );
};

export const RenameListItem = ({ filePath, refresh, setRefresh }) => {
  const [modalOpen, setModalOpen] = useState(false);
  const [success, setSuccess] = useState(false);

  return (
    <>
      <li onClick={() => setModalOpen(true)}>Rename</li>
      <RenameModal
        isOpen={modalOpen}
        close={() => setModalOpen(false)}
        filePath={filePath}
        refresh={refresh}
        setRefresh={setRefresh}
        setSuccess={setSuccess}
      />
      {success && (
        <Toast text="File Renamed Successfully" isVisible={success} />
      )}
    </>
  );
};
