import type { ApiError } from "@/lib/api";

export type FieldErrors = Record<string, string>;

const messages: Record<string, string> = {
  INVALID_INPUT: "Revise os dados informados e tente novamente.",
  UNAUTHORIZED: "E-mail ou senha não conferem. Verifique e tente novamente.",
  FORBIDDEN: "Você não tem permissão para concluir esta ação.",
  NOT_FOUND: "Não encontramos o item solicitado.",
  ALREADY_EXISTS: "Já existe um cadastro com estes dados.",
  CONFLICT: "Não foi possível concluir porque estes dados já estão em uso.",
  FAILED_PRECOND: "Esta ação não pode ser concluída neste momento.",
  RESOURCE_EXHAUST: "Muitas tentativas. Aguarde um instante e tente novamente.",
  TIMEOUT: "A solicitação demorou demais. Tente novamente.",
  UNAVAILABLE: "O serviço está indisponível no momento. Tente novamente mais tarde.",
  INTERNAL: "Ocorreu um erro inesperado. Tente novamente.",
  NETWORK: "Não foi possível conectar. Verifique sua conexão e tente novamente.",
};

export function apiErrorMessage(
  error?: ApiError,
  fallback = "Não foi possível concluir esta ação.",
) {
  return error ? (messages[error.code] ?? fallback) : fallback;
}
export function apiFieldErrors(error?: ApiError): FieldErrors {
  if (error?.code !== "INVALID_INPUT") return {};
  return Object.fromEntries(
    (error.details?.fields ?? []).map((field) => [field, "Revise este campo."]),
  );
}
export function requiredFields(
  values: Record<string, string>,
  labels: Record<string, string>,
): FieldErrors {
  const errors: FieldErrors = {};
  Object.entries(values).forEach(([field, value]) => {
    if (!value.trim()) errors[field] = `${labels[field]} é obrigatório.`;
  });
  return errors;
}
export function invalidEmail(email: string) {
  return !/^\S+@\S+\.\S+$/.test(email);
}
