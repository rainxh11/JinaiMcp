import { join } from 'https://deno.land/std@0.208.0/path/mod.ts';

const STORAGE_DIR = '/app/local-storage';
const SCREENSHOTS_DIR = join(STORAGE_DIR, 'instant-screenshots');

export async function ensureStorageDir(): Promise<void> {
  try {
    await Deno.mkdir(SCREENSHOTS_DIR, { recursive: true });
  } catch {
    // Directory might already exist
  }
}

export async function saveScreenshot(filename: string, data: Buffer): Promise<string> {
  await ensureStorageDir();
  const filepath = join(SCREENSHOTS_DIR, filename);
  await Deno.writeFile(filepath, data);
  return filepath;
}

export async function fileExists(filepath: string): Promise<boolean> {
  try {
    await Deno.stat(filepath);
    return true;
  } catch {
    return false;
  }
}
