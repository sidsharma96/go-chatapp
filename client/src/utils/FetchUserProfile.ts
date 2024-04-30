import { UserProfileUrl } from '../constants/UrlConstants';
import { CredentialError } from './CredentialValidator';

type ProfileData = {
  id: string;
  username: string;
  createdAt: string;
  error: string;
};

export default async function FetchUserProfile() {
  var profile: ProfileData = {
    id: '',
    username: '',
    createdAt: '',
    error: '',
  };
  try {
    const response = await fetch(UserProfileUrl, {
      method: 'GET',
      headers: {
        'content-type': 'application/json',
      },
      credentials: 'include',
    });
    const data = await response.json();
    if (response.ok) {
      profile.id = data.id as string;
      profile.createdAt = data.created_at as string;
      profile.username = data.username as string;
      return profile;
    } else if (response.status === 401) {
      profile.error = CredentialError.unauthenticatedError;
      return profile;
    }
  } catch (err) {
    console.log(err);
    profile.error = CredentialError.catchAllError;
    return profile;
  }
}
