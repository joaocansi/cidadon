"use client";

import { Loader2Icon, MessageCircleReplyIcon, SendIcon, Trash2Icon, XIcon } from "lucide-react";
import { type FormEvent, useMemo, useState } from "react";

import { ImageAttachmentPicker } from "@/components/shared/image-attachment-picker";
import { showToast } from "@/components/shared/toast";
import { UserAvatar } from "@/components/shared/user-avatar";
import { Button } from "@/components/ui/button";
import { DemandImageGallery } from "@/features/demands/components/demand-image-gallery";
import { apiCommentDemand, apiDeleteComment, type DemandComment, type UserRole } from "@/lib/api";
import { apiErrorMessage } from "@/lib/forms";
import { cn } from "@/lib/utils";

type DemandCommentsProps = {
  demandID: number;
  comments: DemandComment[];
  currentUserID?: number;
  userRole?: UserRole;
  officeImageUrl?: string;
  onPublished: () => Promise<void>;
};

export function DemandComments({
  demandID,
  comments,
  currentUserID,
  userRole,
  officeImageUrl,
  onPublished,
}: DemandCommentsProps) {
  const [replyTarget, setReplyTarget] = useState<DemandComment>();
  const roots = useMemo(() => comments.filter((comment) => !comment.parent_id), [comments]);
  const repliesByParent = useMemo(() => {
    const replies = new Map<number, DemandComment[]>();
    comments.forEach((comment) => {
      if (!comment.parent_id) return;
      replies.set(comment.parent_id, [...(replies.get(comment.parent_id) ?? []), comment]);
    });
    return replies;
  }, [comments]);

  return (
    <section className="space-y-5" aria-label="Comentários">
      <CommentComposer demandID={demandID} userRole={userRole} onPublished={onPublished} />
      {roots.length ? (
        <div className="space-y-4">
          {roots.map((comment) => (
            <article key={comment.id} className="min-w-0 space-y-3">
              <CommentCard
                comment={comment}
                officeImageUrl={officeImageUrl}
                currentUserID={currentUserID}
                onDeleted={onPublished}
                onReply={() => setReplyTarget(comment)}
              />
              {replyTarget?.id === comment.id ? (
                <div className="border-lime/70 ml-5 min-w-0 border-l pl-5 sm:ml-10 sm:pl-6">
                  <CommentComposer
                    demandID={demandID}
                    userRole={userRole}
                    replyTo={comment}
                    inline
                    onCancelReply={() => setReplyTarget(undefined)}
                    onPublished={async () => {
                      setReplyTarget(undefined);
                      await onPublished();
                    }}
                  />
                </div>
              ) : null}
              {(repliesByParent.get(comment.id) ?? []).length ? (
                <div className="border-lime/70 ml-5 min-w-0 space-y-3 border-l pl-5 sm:ml-10 sm:pl-6">
                  <p className="text-ink-faint pt-1 text-xs font-semibold tracking-wide uppercase">
                    Respostas
                  </p>
                  {(repliesByParent.get(comment.id) ?? []).map((reply) => (
                    <div key={reply.id} className="min-w-0 space-y-3">
                      <CommentCard
                        comment={reply}
                        officeImageUrl={officeImageUrl}
                        currentUserID={currentUserID}
                        onDeleted={onPublished}
                        onReply={() => setReplyTarget(reply)}
                      />
                      {replyTarget?.id === reply.id ? (
                        <div className="border-lime/70 ml-5 border-l pl-5 sm:ml-10 sm:pl-6">
                          <CommentComposer
                            demandID={demandID}
                            userRole={userRole}
                            replyTo={reply}
                            inline
                            onCancelReply={() => setReplyTarget(undefined)}
                            onPublished={async () => {
                              setReplyTarget(undefined);
                              await onPublished();
                            }}
                          />
                        </div>
                      ) : null}
                    </div>
                  ))}
                </div>
              ) : null}
            </article>
          ))}
        </div>
      ) : (
        <p className="text-ink-soft rounded-xl border border-dashed px-4 py-8 text-center text-sm">
          Ainda não há comentários. Seja a primeira pessoa a contribuir.
        </p>
      )}
    </section>
  );
}

function CommentCard({
  comment,
  officeImageUrl,
  currentUserID,
  onDeleted,
  onReply,
}: {
  comment: DemandComment;
  officeImageUrl?: string;
  currentUserID?: number;
  onDeleted: () => Promise<void>;
  onReply: () => void;
}) {
  const fromOffice =
    comment.author_role === "councillor" || comment.author_role === "office_member";
  const images = comment.images ?? [];
  const [deleting, setDeleting] = useState(false);

  async function remove() {
    if (
      !window.confirm(
        "Excluir este comentário e todas as suas respostas? Esta ação não pode ser desfeita.",
      )
    )
      return;
    setDeleting(true);
    const result = await apiDeleteComment(comment.id);
    setDeleting(false);
    if (!result.ok) {
      showToast(apiErrorMessage(result.error, "Não foi possível excluir o comentário."), "error");
      return;
    }
    showToast("Comentário excluído com sucesso.");
    await onDeleted();
  }
  return (
    <div
      className={cn(
        "rounded-xl border p-4 sm:p-5",
        fromOffice ? "border-lime bg-lime-pale/60 border-l-4" : "border-line bg-card",
      )}
    >
      <div className="flex gap-3">
        <UserAvatar
          name={comment.author_name}
          imageUrl={comment.author_image_url || (fromOffice ? officeImageUrl : undefined)}
          className={cn("size-10", fromOffice && "ring-lime ring-2 ring-offset-2")}
        />
        <div className="min-w-0 flex-1">
          <div className="flex flex-wrap items-center gap-x-2 gap-y-1">
            <p className="text-sm font-semibold">{comment.author_name}</p>
            {fromOffice ? (
              <span className="bg-pine text-card rounded-full px-2 py-0.5 text-[11px] font-bold tracking-wide uppercase">
                Mensagem do gabinete
              </span>
            ) : null}
            <time className="text-ink-faint text-xs" dateTime={comment.created_at}>
              {new Date(comment.created_at).toLocaleString("pt-BR")}
            </time>
          </div>
          {comment.hidden ? (
            <p className="text-ink-soft mt-3 text-sm italic">Conteúdo moderado.</p>
          ) : (
            <>
              {comment.body ? (
                <p className="text-ink-soft mt-3 text-sm leading-6 whitespace-pre-wrap">
                  {comment.body}
                </p>
              ) : null}
              {images.length ? (
                <DemandImageGallery
                  images={images}
                  altPrefix="Imagem anexada ao comentário"
                  className="mt-3"
                />
              ) : null}
            </>
          )}
          <div className="mt-3 flex items-center gap-1">
            <Button variant="ghost" size="sm" className="-ml-2" onClick={onReply}>
              <MessageCircleReplyIcon />
              Responder
            </Button>
            {currentUserID === comment.author_id ? (
              <Button
                variant="ghost"
                size="sm"
                className="text-destructive hover:text-destructive"
                disabled={deleting}
                onClick={() => void remove()}
              >
                {deleting ? <Loader2Icon className="animate-spin" /> : <Trash2Icon />}
                Excluir
              </Button>
            ) : null}
          </div>
        </div>
      </div>
    </div>
  );
}

function CommentComposer({
  demandID,
  userRole,
  replyTo,
  onCancelReply,
  onPublished,
  inline = false,
}: {
  demandID: number;
  userRole?: UserRole;
  replyTo?: DemandComment;
  onCancelReply?: () => void;
  onPublished: () => Promise<void>;
  inline?: boolean;
}) {
  const [body, setBody] = useState("");
  const [images, setImages] = useState<File[]>([]);
  const [saving, setSaving] = useState(false);

  async function submit(event: FormEvent) {
    event.preventDefault();
    if (!body.trim() && !images.length) return;
    setSaving(true);
    const result = await apiCommentDemand(demandID, {
      body,
      images,
      parent_id: replyTo?.id,
    });
    setSaving(false);
    if (!result.ok) {
      showToast(apiErrorMessage(result.error, "Não foi possível publicar o comentário."), "error");
      return;
    }
    setBody("");
    setImages([]);
    showToast(replyTo ? "Resposta publicada." : "Comentário publicado.");
    await onPublished();
  }

  const fromOffice = userRole === "councillor" || userRole === "office_member";
  return (
    <form
      onSubmit={submit}
      className={
        inline ? "bg-paper-2 rounded-xl p-4" : "border-line bg-paper-2 rounded-xl border p-4 sm:p-5"
      }
    >
      <div className="mb-3 flex items-center justify-between gap-3">
        <p className="text-sm font-semibold">
          {replyTo ? `Respondendo a ${replyTo.author_name}` : "Adicionar comentário"}
        </p>
        {replyTo ? (
          <Button type="button" variant="ghost" size="sm" onClick={onCancelReply}>
            <XIcon /> Cancelar resposta
          </Button>
        ) : null}
      </div>
      {fromOffice ? (
        <p className="text-pine mb-3 text-xs font-semibold">
          Sua contribuição será identificada como uma mensagem do gabinete.
        </p>
      ) : null}
      <textarea
        value={body}
        onChange={(event) => setBody(event.target.value)}
        placeholder={
          replyTo
            ? "Escreva sua resposta"
            : "Compartilhe uma atualização, dúvida ou informação útil"
        }
        className="field-textarea min-h-28"
        maxLength={2000}
      />
      <div className="mt-4 flex flex-wrap items-center justify-between gap-3">
        <ImageAttachmentPicker files={images} onChange={setImages} disabled={saving} />
        <Button type="submit" disabled={saving || (!body.trim() && !images.length)}>
          {saving ? <Loader2Icon className="animate-spin" /> : <SendIcon />}
          Publicar
        </Button>
      </div>
    </form>
  );
}
