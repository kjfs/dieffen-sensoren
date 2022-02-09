function Prompt() {
  let toast = function (c) {
    const { msg = "", icon = "success", position = "top-end" } = c;
    const Toast = Swal.mixin({
      toast: true,
      title: msg,
      icon: icon,
      position: position,
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
    const {
      msg = "",
      icon = "success",
      position = "top-end",
      title = "",
      footer = "",
    } = c;
    Swal.fire({
      icon: icon,
      title: title,
      text: msg,
      footer: footer,
    });
  };
  let error = function (c) {
    const {
      msg = "",
      icon = "error",
      position = "top-end",
      title = "",
      footer = "",
    } = c;
    Swal.fire({
      icon: icon,
      title: title,
      text: msg,
      footer: footer,
    });
  };

  async function custom(c) {
    const { icon = "", msg = "", title = "", showConfirmButton = true } = c;

    const { value: result } = await Swal.fire({
      //title: 'Multiple inputs',
      title: title,
      // html:
      // '<input id="swal-input1" class="swal2-input">' +
      // '<input id="swal-input2" class="swal2-input">',
      html: msg,
      backdrop: false,
      focusConfirm: false,
      showCancelButton: true,
      icon: icon,
      showConfirmButton: showConfirmButton,
      willOpen: () => {
        // const elem = document.getElementById("reservation-dates-modal");
        // const rp = new DateRangePicker(elem, {
        //     format: "yyyy/mm/dd",
        //     minDate: 1,
        //     todayBtn: true,
        //     todayHighlight: true,
        //     weekStart: 1,
        //     showOnFocus: true,
        // });
        if (c.willOpen !== undefined) {
          c.willOpen();
        }
      },
      // preConfirm: () => {
      //   return [
      //     document.getElementById("start").value,
      //     document.getElementById("end").value,
      //   ];
      // },
      didOpen: () => {
        // return [
        //     document.getElementById("start").removeAttribute("disabled"),
        //     document.getElementById("end").removeAttribute("disabled")
        // ]
        if (c.didOpen !== undefined) {
          c.didOpen();
        }
      },
    });

    // if (formValues) {
    //     Swal.fire(JSON.stringify(formValues))
    // }
    if (result) {
      if (result.dismiss !== Swal.DismissReason.cancel) {
        if (result.value !== "") {
          if (c.callback !== undefined) {
            c.callback(result);
          }
        } else {
          c.callback(false);
        }
      } else {
        c.callback(false);
      }
    }
  }

  return {
    toast: toast,
    success: success,
    error: error,
    custom: custom,
  };
}

function lokko(token, sensorID) {
  document
    .getElementById("register-yourself-button")
    .addEventListener("click", function () {
      let html = `
        <form id="check-availability-form" action="" method="post" novalidate class="needs-valiation">
            <div class="form-row">
                <div class="col">
                    <div class="form-row" id="reservation-dates-modal">
                        <div class="col">
                            <input disabled required class="form-control" type="text" name="start" id="start" placeholder="StartDate">
                        </div>
                        <div class="col">
                            <input disabled required class="form-control" type="text" name="end" id="end" placeholder="EndDate">
                        </div>
                    </div>
                </div>
            </div>
        </form>
        `;

      // attention.error({msg: "Das ist eine error", title: "error", footer: "Footer sind nice"});
      // attention.custom({msg: html, title: "Choose your dates"});
      attention.custom({
        msg: html,
        title: "Choose your dates",

        willOpen: () => {
          const elem = document.getElementById("reservation-dates-modal");
          const rp = new DateRangePicker(elem, {
            format: "yyyy-mm-dd",
            minDate: new Date(),
            todayBtn: true,
            todayHighlight: true,
            weekStart: 1,
            showOnFocus: true,
          });
        },

        didOpen: () => {
          return [
            document.getElementById("start").removeAttribute("disabled"),
            document.getElementById("end").removeAttribute("disabled"),
          ];
        },

        callback: function (result) {
          console.log("called");

          let form = document.getElementById("check-availability-form");
          let formData = new FormData(form);
          formData.append("csrf_token", token);
          formData.append("sensor_id", sensorID);

          fetch("/search-json", {
            method: "post",
            body: formData,
          })
            .then((response) => response.json())
            .then((data) => {
              if (data.ok) {
                attention.custom({
                  icon: "success",
                  title: "Data exist",
                  showConfirmButton: false,
                  msg:
                    "<p> Data is available </p>" +
                    '<p><a href="/book-room?id=' +
                    data.sensor_id +
                    "&s=" +
                    data.start_date +
                    "&e=" +
                    data.end_date +
                    '" class="btn btn-primary">' +
                    "Data exist! </a></p>",
                });
              } else {
                console.log("No Availability / No Data");
                attention.error({
                  msg: "No Availability / No Data",
                });
              }
            });
        },
      });
    });
}
