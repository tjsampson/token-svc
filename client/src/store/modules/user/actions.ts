import { ActionTree } from "vuex";
import { UserState } from "./types";
import { RootState } from "@/store/root";
import { userService } from "@/_services";
import { router } from "@/router";

const alertActionRoot = "alert/error";

export const actions: ActionTree<UserState, RootState> = {
  async login({ dispatch, commit }, { email, password }) {
    try {
      commit("LOGIN_REQUEST", { email });
      const user = await userService.login(email, password);
      commit("LOGIN_SUCCESS", user);
      router.push("/");
    } catch (error) {
      commit("LOGIN_FAILURE", error);
      dispatch(alertActionRoot, error, { root: true });
    }
  },
  logout({ dispatch, commit }) {
    commit("LOGOUT_REQUEST");
    userService.logout();
    commit("LOGOUT_SUCCESS");
  },
  async register({ dispatch, commit }, { email, password, confirm_password }) {
    try {
      commit("REGISTER_REQUEST", { email });
      const user = await userService.register(
        email,
        password,
        confirm_password
      );
      commit("REGISTER_SUCCESS", user);
      router.push("/home");
    } catch (error) {
      commit("REGISTER_FAILURE", error);
      dispatch(alertActionRoot, error, { root: true });
    }
  },
  async list({ dispatch, commit }) {
    commit("LIST_USERS_REQUEST");
    try {
      const users = await userService.list();
      commit("LIST_USERS_SUCCESS", users);
    } catch (error) {
      commit("LIST_USERS_FAILURE", error);
      dispatch(alertActionRoot, error, { root: true });
    }
  }
};
