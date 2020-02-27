export function authHeader() {
  // return authorization header with jwt token
  let user = JSON.parse(localStorage.getItem("user") || "{}");

  if (user && user.token) {
    return { Authorization: "Bearer " + user.token };
  } else {
    return {};
  }
}

export function userToken(key: string): string {
  // return authorization header with jwt token
  let accessToken = JSON.parse(localStorage.getItem(key) || "{}");

  if (accessToken) {
    return accessToken;
  } else {
    return "";
  }
}
