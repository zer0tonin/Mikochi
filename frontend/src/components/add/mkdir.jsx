import { h } from "preact";
import { useContext, useState } from "preact/hooks";
import {signal} from "@preact/signals";

import { AuthContext } from "../../jwt";
import Modal, { ModalContent, ModalHeader } from "../modal";
import "./style.css";
import handleError from "../../error";
import Toast from "../toast";

export const mkdirOpen = signal(false);

const MkdirModal = ({
  isOpen,
  close,
  dirPath,
  refresh,
  setRefresh,
  setSuccess,
}) => {
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
        },
      );

      if (response.status !== 200) {
        setError(await handleError(response));
        return;
      }
      setError("");
      setName("");
      setRefresh(refresh + 1);
      close();
      setSuccess(true);
      setTimeout(() => {
        setSuccess(false);
      }, 2000);
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
            class="textInput"
            onChange={(e) => setName(e.target.value)}
          />
          <button type="submit" class="submit">
            Create
          </button>
          {error !== "" && <div class="error">{error}</div>}
        </form>
      </ModalContent>
    </Modal>
  );
};

const Mkdir = ({ dirPath, refresh, setRefresh }) => {
  const [success, setSuccess] = useState(false); //into signal

  return (
    <>
      <MkdirModal
        isOpen={mkdirOpen.value}
        close={() => mkdirOpen.value = false}
        dirPath={dirPath}
        refresh={refresh}
        setRefresh={setRefresh}
        setSuccess={setSuccess}
      />
      {success && (
        <Toast text="Directory Created Successfully" isVisible={success} />
      )}
    </>
  );
};

export default Mkdir;
