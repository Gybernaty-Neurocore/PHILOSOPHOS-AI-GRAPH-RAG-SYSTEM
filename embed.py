import sys
import json
import traceback
from sentence_transformers import SentenceTransformer

try:
    model = SentenceTransformer('all-mpnet-base-v2')
except Exception as e:
    print(f"Error loading model: {str(e)}", file=sys.stderr)
    traceback.print_exc(file=sys.stderr)
    sys.exit(1)

def get_embedding(text):
    try:
        embedding = model.encode(text, convert_to_numpy=True).tolist()
        return embedding
    except Exception as e:
        print(f"Error encoding text: {str(e)}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    try:
        text = sys.argv[1].encode('utf-8').decode('utf-8')
        embedding = get_embedding(text)
        print(json.dumps(embedding))
    except Exception as e:
        print(f"Error: {str(e)}", file=sys.stderr)
        traceback.print_exc(file=sys.stderr)
        sys.exit(1)