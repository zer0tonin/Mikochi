import { h } from "preact";
import { useContext, useState } from "preact/hooks";

import { AuthContext } from "../../utils/jwt";
import style from "./style.css";
import { BigIcon } from "../icon";
import Modal, { ModalContent, ModalHeader } from "../modal";

const MkdirModal = ({ isOpen, close, dirPath, refresh, setRefresh }) => {
  const { jwt } = useContext(AuthContext);
  const [error, setError] = useState("");
  const [name, setName] = useState("");

  if (!isOpen) {
    return null;
  }

  const onSubmit = (e) => {
    e.preventDefault();
    const putMkdir = async () => {
      const response = await fetch(
        dirPath != "" ? `/api/mkdir/${dirPath}/${name}` : `/api/mkdir/${name}`,
        {
          method: "PUT",
          headers: {
            Accept: "application/json",
            Authorization: `Bearer ${jwt}`,
          },
        }
      );

      if (response.status !== 200) {
        setError(response.json()["err"]);
        return;
      }
      setError("");
      setName("");
      setRefresh(refresh + 1);
      close();
    };
    putMkdir();
  };

  return (
    <Modal isOpen={isOpen} close={close}>
      <ModalHeader close={close}>Create directory</ModalHeader>
      <ModalContent>
        <form onSubmit={onSubmit}>
          <input
            type="text"
            value={name}
            class={style.textInput}
            onChange={(e) => setName(e.target.value)}
          />
          <button type="submit" class={style.submit}>
            Create
          </button>
          {error !== "" && <div class={style.error}>{error}</div>}
        </form>
      </ModalContent>
    </Modal>
  );
};

const Mkdir = ({ dirPath, refresh, setRefresh }) => {
  const [modalOpen, setModalOpen] = useState(false);
  return (
    <>
      <div class={style.floatingMkdir}>
        <BigIcon
          name="folder-add"
          onClick={() => setModalOpen(true)}
          title="Create directory"
        />
      </div>
      <MkdirModal
        isOpen={modalOpen}
        close={() => setModalOpen(false)}
        dirPath={dirPath}
        refresh={refresh}
        setRefresh={setRefresh}
      />
    </>
  );
};

export default Mkdir;
