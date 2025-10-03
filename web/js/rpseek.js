(async function () {
  const res = await fetch("/search", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      query: "3P_0 AFTER AND AVAILABLE BARYONS BEEN CALCULATED CHARMED CHEN",
    }),
  });

  console.log(res);
})();
