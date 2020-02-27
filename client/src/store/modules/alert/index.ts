import { Module } from "vuex";
import { getters } from "./getters";
import { actions } from "./actions";
import { mutations } from "./mutations";
import { AlertState } from "./types";
import { RootState } from "@/store/root";

export const state: AlertState = {
  type: null,
  message: null
};

const namespaced: boolean = true;

export const user: Module<AlertState, RootState> = {
  namespaced,
  state,
  getters,
  actions,
  mutations
};

export default user;
