import { X } from "lucide-react";
import { ReactElement, ReactEventHandler, useEffect, useRef } from "react";

function Modal({ openModal, closeModal, children }: { openModal: boolean, closeModal: ReactEventHandler, children: ReactElement }) {
  const ref = useRef<HTMLDialogElement>(null as any as HTMLDialogElement);

  useEffect(() => {
    if (openModal) {
      ref.current?.showModal();
    } else {
      ref.current?.close();
    }
  }, [openModal]);

  return (
    <dialog
      ref={ref}
      onCancel={closeModal}
      style={{ position: "relative" }}
    >
      {children}
      <button onClick={closeModal} className="close-button">
        <X />
      </button>
    </dialog>
  );
}

export default Modal;
