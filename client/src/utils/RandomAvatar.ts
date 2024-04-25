const avatarSeeds: string[] = [
  'Abby',
  'Sophie',
  'Snowball',
  'Casper',
  'Tigger',
  'Molly',
  'Sasha',
  'Chloe',
  'Boots',
  'Missy',
  'Patches',
  'Oreo',
  'Willow',
  'Toby',
  'Mia',
  'Nala',
  'Gizmo',
  'Bella',
  'Spooky',
  'Muffin',
];

export default function FetchRandomAvatarImage(): string {
  const baseUrl: string = 'https://api.dicebear.com/8.x/avataaars-neutral/svg?';
  const imageSrc: string =
    baseUrl +
    'seed=' +
    avatarSeeds[Math.floor(Math.random() * avatarSeeds.length)] +
    '&size=275';
  return imageSrc;
}
