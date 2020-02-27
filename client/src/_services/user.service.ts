import config from "../config";
import { userToken } from "../_helpers";

const userAccessTokenStorageKey = "user:access";
const userRefreshTokenStorageKey = "user:refresh";
const userStorageKey = "user";

export const userService = {
  login,
  logout,
  list,
  register,
  userAccessTokenStorageKey,
  userRefreshTokenStorageKey,
  userStorageKey
};

async function register(
  email: string,
  password: string,
  confirm_password: string
) {
  const requestOptions = {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password, confirm_password })
  };

  const response = await fetch(
    `${config.BASE_API_URL}/register`,
    requestOptions
  );

  const registerResponse = await handleResponse(response);

  if (registerResponse.id > 0) {
    localStorage.setItem("user", JSON.stringify(registerResponse));
  }

  return registerResponse;
}

async function login(email: string, password: string) {
  const requestOptions = {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password })
  };

  const response = await fetch(`${config.BASE_API_URL}/login`, requestOptions);
  const loginResponse = await handleResponse(response);
  // login successful if there's a jwt access_token in the response
  if (loginResponse.access_token) {
    // store jwt access_token in local storage to keep user logged in between page refreshes
    localStorage.setItem(
      userAccessTokenStorageKey,
      JSON.stringify(loginResponse.access_token)
    );
  }
  if (loginResponse.refresh_token) {
    // store jwt refresh_token in local storage to give client ability to refresh the access token on expiration
    localStorage.setItem(
      userRefreshTokenStorageKey,
      JSON.stringify(loginResponse.refresh_token)
    );
  }
  return loginResponse;
}

function logout() {
  // remove user from local storage to log user out
  localStorage.removeItem(userAccessTokenStorageKey);
  localStorage.removeItem(userRefreshTokenStorageKey);
  localStorage.removeItem(userStorageKey);
}

async function list() {
  const requestOptions = {
    method: "GET",
    headers: { Authorization: "Bearer " + userToken(userAccessTokenStorageKey) }
  };

  const response = await fetch(`${config.BASE_API_URL}/users`, requestOptions);
  const listResponse = await handleResponse(response);

  return listResponse;
}

async function handleResponse(response: Response) {
  const text = await response.text();
  const data = text && JSON.parse(text);

  if (!response.ok) {
    if (response.status === 401) {
      // auto logout if 401 response returned from api
      logout();
      // location.reload(true);
    }

    if (data && data.messages && data.messages.length > 0) {
      return Promise.reject(data.messages);
    } else {
      const error = (data && data.message) || response.statusText;
      return Promise.reject(error);
    }
  }
  return data;
}
