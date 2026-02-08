export const checkAndGetUserCache = () => {
  const cached = localStorage.getItem("userData");
  if (!cached) return null;

  const user = JSON.parse(cached);
  if (
    user.userId &&
    user.login &&
    user.token &&
    user.avatarLink !== undefined
  ) {
    return {
      ...user,
      valid: true,
    };
  }

  return null;
};

export const saveUserToLocalStorage = (state) => {
  localStorage.setItem("userData", JSON.stringify(state));
};

export const isUserValid = (state) => state.user.valid;