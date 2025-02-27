function Prompt() {
    let toast = function (c) {
      const { msg = "", icon = "success", position = "top-end" } = c;

      const Toast = Swal.mixin({
        toast: true,
        title: msg,
        position: position,
        icon: icon,
        showConfirmButton: false,
        timer: 3000,
        timerProgressBar: true,
        didOpen: (toast) => {
          toast.addEventListener("mouseenter", Swal.stopTimer);
          toast.addEventListener("mouseleave", Swal.resumeTimer);
        },
      });

      Toast.fire({});
    };

    let success = function (c) {
      const { msg = "", title = "", footer = "" } = c;

      Swal.fire({
        icon: "success",
        title: title,
        text: msg,
        footer: footer,
      });
    };

    let error = function (c) {
      const { msg = "", title = "", footer = "" } = c;

      Swal.fire({
        icon: "error",
        title: title,
        text: msg,
        footer: footer,
      });
    };

    async function custom(c) {
      const { icon = "", msg = "", title = "", showConfirmButton = true} = c;

      const { value: result } = await Swal.fire({
        icon: icon,
        title: title,
        html: msg,
        backdrop: false,
        focusConfirm: false,
        showConfirmButton: showConfirmButton,
        showCancelButton: true,
        willOpen: () => {
          if (c.willOpen) {
            c.willOpen();
          }
        },
        didOpen: () => {
          if (c.didOpen) {
            c.didOpen();
          }
        },
      });
      console.log(result, 44)
      if (!result) {
        c.callback(false);
      } else if (typeof result === "boolean") {
        c.callback(true)
      } else if (result.dismiss === Swal.DismissReason.cancel) {
        c.callback(false);
      } else if (!result.value) {
        c.callback(false);
      } else c.callback(result);
    }

    return {
      toast: toast,
      success: success,
      error: error,
      custom: custom,
    };
  }