import { GetterTree } from "vuex";
import { UserState } from "./types";
import { RootState } from "@/store/root";

export const getters: GetterTree<UserState, RootState> = {
  USERS: state => state.users
};
