import { createSlice } from '@reduxjs/toolkit';
import { checkAndGetUserCache, saveUserToLocalStorage } from '../localStorage';

const cachedUser = checkAndGetUserCache();

const initialState = cachedUser || {
  userId: null,
  login: '',
  phoneNumber: '',
  avatarLink: null,
  searchPrivacy: false,
  token: '',
  valid: false,
};

const meSlice = createSlice({
  name: 'user',
  initialState,
  reducers: {
    setUser: (state, action) => {
      const { userId, login, phoneNumber, avatarLink, searchPrivacy } = action.payload;
      state.userId = userId;
      state.login = login;
      state.phoneNumber = phoneNumber;
      state.avatarLink = avatarLink;
      state.searchPrivacy = searchPrivacy;
      state.valid = true;
      saveUserToLocalStorage(state);
    },
    setUserInAuth: (state, action) => {
      const data = action.payload;
      state.userId = data.user_id;
      state.login = data.login;
      state.token = data.token;
      state.valid = true;

      saveUserToLocalStorage(state);
    },
    updateUser: (state, action) => {
      Object.assign(state, action.payload);
      state.valid = true;

      saveUserToLocalStorage(state);
    },
  },
});

export const { setUser, updateUser, setUserInAuth } = meSlice.actions;
export default meSlice.reducer;
