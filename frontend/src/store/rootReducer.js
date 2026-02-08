import { combineReducers } from '@reduxjs/toolkit';
import chatsReducer from './slices/chatSlice';
import meSliceReducer from './slices/meSlice';

const appReducer = combineReducers({
  chats: chatsReducer,
  user: meSliceReducer,
});

const rootReducer = (state, action) => {
  if (action.type === 'LOGOUT') {
    state = undefined;
  }

  return appReducer(state, action);
};

export default rootReducer;
