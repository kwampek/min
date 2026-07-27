import { createSlice } from "@reduxjs/toolkit";

const initialState = {
    loading: {
        users: false,
    },

    chat: null,
    users: [],
};

const chatDialogSlice = createSlice({
    name: "chatDialog",

    initialState,

    reducers: {
        openDialog(state, action) {
            const { mode, type, chat } = action.payload;

            state.chat = chat ?? null;

            state.users = [];
            state.loading.users = true;
        },

        closeDialog() {
            return initialState;
        },

        setUsers(state, action) {
            state.users = action.payload;
            state.loading.users = false;
        },
    },
});

export const {
    openDialog,
    closeDialog,
    setUsers,
} = chatDialogSlice.actions;

export default chatDialogSlice.reducer;