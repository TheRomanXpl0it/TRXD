"use client"
import {
  useState
} from "react"
import {
  toast
} from "sonner"
import {
  useForm
} from "react-hook-form"
import {
  zodResolver
} from "@hookform/resolvers/zod"
import * as z from "zod"
import {
  cn
} from "@/lib/utils"
import {
  Button
} from "@/components/ui/button"
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import {
  Input
} from "@/components/ui/input"
import {
  Textarea
} from "@/components/ui/textarea"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from "@/components/ui/select"
import {
  TagsInput
} from "@/components/ui/tags-input"
import {
  Switch
} from "@/components/ui/switch"
import {
  CloudUpload,
  Paperclip
} from "lucide-react"
import {
  FileInput,
  FileUploader,
  FileUploaderContent,
  FileUploaderItem
} from "@/components/ui/file-upload"
import {
  MultiSelector,
  MultiSelectorContent,
  MultiSelectorInput,
  MultiSelectorItem,
  MultiSelectorList,
  MultiSelectorTrigger
} from "@/components/ui/multi-select"
import { ChallengeProps } from "./challenge"

const formSchema = z.object({
  title: z.string().min(1),
  flag: z.string().min(1),
  description: z.string(),
  difficulty: z.string(),
  tags: z.array(z.string()).nonempty("Please at least one item"),
  hidden: z.boolean(),
  files: z.string(),
  authors: z.array(z.string())
});

function displayAuthors(authors: string[]){
  return (
    authors.map((author, index) => (
      <MultiSelectorItem key={index} value={author}>{author}</MultiSelectorItem>
    ))
  )
}

export function ChallengeForm( { challengeProp, auth } : { 
  challengeProp?: ChallengeProps,
  auth: {
    username: string;
    password: string;
    accessToken: string;
    roles: string[];
}}) {

  let defaultAuthors: string[] | undefined = [auth.username];
  let defaultTitle: string | undefined = "";
  let defaultDescription: string | undefined = "";
  let defaultDifficulty: string | undefined = "Easy";
  let defaultTags: string[] | undefined = ["Easy"];
  let defaultHidden: boolean | undefined = true;
  let defaultFlag : string | undefined = "";
  let authors: string[] = [];


  if ( challengeProp ){
    console.log(challengeProp);
    challengeProp.challenge.authors ? defaultAuthors = challengeProp.challenge.authors : [];
    challengeProp.challenge.description ? defaultDescription = challengeProp.challenge.description : "";
    challengeProp.challenge.difficulty ? defaultDifficulty = challengeProp.challenge.difficulty : "Easy";
    challengeProp.challenge.tags ? defaultTags = challengeProp.challenge.tags : ["Easy"];
    challengeProp.challenge.hidden!==undefined ? defaultHidden = challengeProp.challenge.hidden : true;
    defaultTitle = challengeProp.challenge.title;
    defaultFlag = challengeProp.challenge.flag;
  }


  const [files, setFiles] = useState < File[] | null > (null);

  const dropZoneConfig = {
    maxFiles: 5,
    maxSize: 1024 * 1024 * 4,
    multiple: true,
  };
  const form = useForm < z.infer < typeof formSchema >> ({
    resolver: zodResolver(formSchema),
    defaultValues: {
      "title" : defaultTitle,
      "flag": defaultFlag,
      "description" : defaultDescription,
      "difficulty" : defaultDifficulty,
      "tags": defaultTags,
      "authors": defaultAuthors,
      "hidden": defaultHidden,
    },
  })

  function onSubmit(values: z.infer < typeof formSchema > ) {
    try {
      console.log(values);
      toast(
        <pre className="mt-2 w-[340px] rounded-md bg-slate-950 p-4">
          <code className="text-white">{JSON.stringify(values, null, 2)}</code>
        </pre>
      );
    } catch (error) {
      console.error("Form submission error", error);
      toast.error("Failed to submit the form. Please try again.");
    }
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8 max-w-3xl mx-auto py-10">
        
        <FormField
          control={form.control}
          name="title"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Title</FormLabel>
              <FormControl>
                <Input 
                placeholder=""
                
                type="text"
                {...field} />
              </FormControl>
              <FormDescription>The title of the challenge</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        
        <FormField
          control={form.control}
          name="flag"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Flag</FormLabel>
              <FormControl>
                <Input 
                placeholder="TRX{...}"
                
                type="text"
                {...field} />
              </FormControl>
              <FormDescription>Correct Flag</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        
        <FormField
          control={form.control}
          name="description"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Description</FormLabel>
              <FormControl>
                <Textarea
                  placeholder=""
                  className="resize-none"
                  {...field}
                />
              </FormControl>
              <FormDescription>The description of the challenge</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        
        <FormField
          control={form.control}
          name="difficulty"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Difficulty</FormLabel>
              <Select onValueChange={field.onChange} defaultValue={field.value}>
                <FormControl>
                  <SelectTrigger>
                    <SelectValue placeholder="" />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  <SelectItem value="Easy">Easy</SelectItem>
                  <SelectItem value="Medium">Medium</SelectItem>
                  <SelectItem value="Hard">Hard</SelectItem>
                  <SelectItem value="Insane">Insane</SelectItem>
                </SelectContent>
              </Select>
                <FormDescription>Difficulty of the challenge</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        
        <FormField
          control={form.control}
          name="tags"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Tags</FormLabel>
              <FormControl>
                <TagsInput
                  value={field.value}
                  onValueChange={field.onChange}
                  placeholder="Enter your tags"
                />
              </FormControl>
              <FormDescription>Add tags to the challenge.</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        
          <FormField
              control={form.control}
              name="hidden"
              render={({ field }) => (
                <FormItem className="flex flex-row items-center justify-between rounded-lg border p-4">
                  <div className="space-y-0.5">
                    <FormLabel>Hidden</FormLabel>
                    <FormDescription>Challenge is hidden, you can edit it later on inside the challenge card.</FormDescription>
                  </div>
                  <FormControl>
                    <Switch
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  </FormControl>
                </FormItem>
              )}
            />
        
            <FormField
              control={form.control}
              name="files"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Challenge files</FormLabel>
                  <FormControl>
                    <FileUploader
                      value={files}
                      onValueChange={setFiles}
                      dropzoneOptions={dropZoneConfig}
                      className="relative bg-background rounded-lg p-2"
                    >
                      <FileInput
                        id="fileInput"
                        className="outline-dashed outline-1 outline-slate-500"
                      >
                        <div className="flex items-center justify-center flex-col p-8 w-full ">
                          <CloudUpload className='text-gray-500 w-10 h-10' />
                          <p className="mb-1 text-sm text-gray-500 dark:text-gray-400">
                            <span className="font-semibold">Click to upload</span>
                            &nbsp; or drag and drop
                          </p>
                          <p className="text-xs text-gray-500 dark:text-gray-400">
                            (.zip)
                          </p>
                        </div>
                      </FileInput>
                      <FileUploaderContent>
                        {files &&
                          files.length > 0 &&
                          files.map((file, i) => (
                            <FileUploaderItem key={i} index={i}>
                              <Paperclip className="h-4 w-4 stroke-current" />
                              <span>{file.name}</span>
                            </FileUploaderItem>
                          ))}
                      </FileUploaderContent>
                    </FileUploader>
                  </FormControl>
                  <FormDescription>Select a file to upload.</FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
        
           <FormField
              control={form.control}
              name="authors"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Authors</FormLabel>
                  <FormControl>
                    <MultiSelector
                      values={field.value}
                      onValuesChange={field.onChange}
                      loop
                      className="max-w-xs"
                    >
                      <MultiSelectorTrigger>
                        <MultiSelectorInput placeholder="Select Authors" />
                      </MultiSelectorTrigger>
                      <MultiSelectorContent>
                      <MultiSelectorList>
                        { displayAuthors(authors) }
                      </MultiSelectorList>
                      </MultiSelectorContent>
                    </MultiSelector>
                  </FormControl>
                  <FormDescription>Select multiple options.</FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
        <Button type="submit">Submit</Button>
      </form>
    </Form>
  )
}