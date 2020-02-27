import { Module } from "vuex";
import { getters } from "./getters";
import { actions } from "./actions";
import { mutations } from "./mutations";
import { UserState } from "./types";
import { RootState } from "@/store/root";

const loggedInUser = JSON.parse(localStorage.getItem("user") || "{}");

export const state: UserState = {
  email: "",
  password: "",
  users: [],
  isLoggedIn: () => {
    return false;
  },
  loading: false,
  record: {},
  authenticated: false
};

const namespaced: boolean = true;

export const user: Module<UserState, RootState> = {
  namespaced,
  state,
  getters,
  actions,
  mutations
};

export default user;
